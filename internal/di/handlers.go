//nolint:ireturn
package di

import (
	"context"
	"electrotech/internal/config"
	"electrotech/internal/email"
	"electrotech/internal/handlers/auth"
	catalogHandlers "electrotech/internal/handlers/catalog"
	v2 "electrotech/internal/handlers/catalog/v2"
	"electrotech/internal/handlers/contact"
	"electrotech/internal/handlers/orders"
	"electrotech/internal/handlers/user"
	"electrotech/internal/repository/catalog"
	ordersRepo "electrotech/internal/repository/orders"
	"electrotech/internal/repository/users"
	"electrotech/internal/server"
	"fmt"

	"charm.land/log/v2"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func newServer(
	lc fx.Lifecycle,
	config *config.Config,
	catalogRepo *catalog.Repo,
	logger *log.Logger,
	contactHandler *contact.ContactUsHandler,
	ordersHandler *orders.Handler,
	userHandler *user.Handler,
	authHandler *auth.Handler,
	catalogHandler *catalogHandlers.Handler,
	catalogV2Handler *v2.Handler,
) *server.HTTPServer {
	srv := server.NewHTTPServer(
		config,
		logger,
		catalogRepo,
		contactHandler,
		ordersHandler,
		userHandler,
		authHandler,
		catalogHandler,
		catalogV2Handler,
	)

	lc.Append(fx.StartHook(func(ctx context.Context) {
		go func() {
			err := srv.Run()
			if err != nil {
				logger.Error("Failed run HTTP server", "error", err)
			}
		}()
	}))

	return srv
}

func provideHandlers() fx.Option {
	return fx.Options(
		fx.Provide(newContactHandler),
		fx.Provide(newOrdersHandler),
		fx.Provide(newUserHandler),
		fx.Provide(newAuthHandler),
		fx.Provide(newCatalogHandler),
		fx.Provide(newCatalogV2Handler),
	)
}

func newContactHandler(emailService *email.Service, logger *log.Logger, db *gorm.DB) *contact.ContactUsHandler {
	return contact.NewContactUsHandler(emailService, logger, db)
}

func newOrdersHandler(
	logger *log.Logger,
	emailService *email.Service,
	catalogRepo *catalog.Repo,
	ordersRepo *ordersRepo.Repo,
	usersRepo *users.Repo,
) *orders.Handler {
	return orders.NewHandler(logger, catalogRepo, ordersRepo, usersRepo, emailService)
}

func newUserHandler(logger *log.Logger, usersRepo *users.Repo) *user.Handler {
	return user.NewHandler(logger, usersRepo)
}

func newAuthHandler(logger *log.Logger, config *config.Config, usersRepo *users.Repo) (*auth.Handler, error) {
	h, err := auth.NewHandler(logger, &config.Auth, usersRepo)
	if err != nil {
		return nil, fmt.Errorf("new auth handler: %w", err)
	}

	return h, nil
}

func newCatalogHandler(repo *catalog.Repo) *catalogHandlers.Handler {
	return catalogHandlers.NewHandler(repo)
}

func newCatalogV2Handler(repo *catalog.Repo) *v2.Handler {
	return v2.NewHandler(repo)
}
