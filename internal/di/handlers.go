//nolint:ireturn
package di

import (
	"context"
	"electrotech/internal/config"
	"electrotech/internal/handlers"
	"electrotech/internal/handlers/auth"
	catalogHandlers "electrotech/internal/handlers/catalog"
	v2 "electrotech/internal/handlers/catalog/v2"
	"electrotech/internal/handlers/contact"
	"electrotech/internal/repository/catalog"
	"electrotech/internal/repository/users"
	"electrotech/internal/server"
	"fmt"

	"electrotech/internal/services/feedback"
	"electrotech/internal/services/health"
	"electrotech/internal/services/order"
	userS "electrotech/internal/services/user"

	"charm.land/log/v2"
	"go.uber.org/fx"
)

func newServer(
	lc fx.Lifecycle,
	config *config.Config,
	catalogRepo *catalog.Repo,
	logger *log.Logger,
	contactHandler *contact.ContactUsHandler,
	authHandler *auth.Handler,
	catalogHandler *catalogHandlers.Handler,
	catalogV2Handler *v2.Handler,
	handler *handlers.Handler,
	middleware *handlers.AuthMiddleware,
) (*server.HTTPServer, error) {
	srv, err := server.NewHTTPServer(
		config,
		logger,
		catalogRepo,
		contactHandler,
		authHandler,
		catalogHandler,
		catalogV2Handler,
		handler,
		middleware,
	)
	if err != nil {
		return nil, fmt.Errorf("new HTTP server: %w", err)
	}

	lc.Append(fx.StartHook(func(ctx context.Context) {
		go func() {
			err := srv.Run()
			if err != nil {
				logger.Error("Failed run HTTP server", "error", err)
			}
		}()
	}))

	return srv, nil
}

func provideHandlers() fx.Option {
	return fx.Options(
		fx.Provide(newContactHandler),
		fx.Provide(newAuthHandler),
		fx.Provide(newCatalogHandler),
		fx.Provide(newCatalogV2Handler),
	)
}

func newContactHandler(logger *log.Logger, feedbackService *feedback.Service) *contact.ContactUsHandler {
	return contact.NewContactUsHandler(logger, feedbackService)
}

func newAuthHandler(
	logger *log.Logger,
	config *config.Config,
	usersRepo *users.Repo,
	userService *userS.Service,
) (*auth.Handler, error) {
	h, err := auth.NewHandler(logger, &config.Auth, usersRepo, userService)
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

func newAuthMiddleware(config *config.Config, service *userS.Service) *handlers.AuthMiddleware {
	return handlers.NewAuthMiddleware(&config.Auth, service)
}

func newHandlers(
	logger *log.Logger,
	orderService *order.Service,
	userService *userS.Service,
	healthService *health.Service,
) *handlers.Handler {
	return handlers.New(logger, orderService, userService, healthService)
}
