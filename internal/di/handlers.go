//nolint:ireturn
package di

import (
	"context"
	"electrotech/internal/config"
	"electrotech/internal/email"
	"electrotech/internal/handlers/contact"
	"electrotech/internal/handlers/orders"
	"electrotech/internal/handlers/user"
	"electrotech/internal/repository/catalog"
	"electrotech/internal/server"
	"electrotech/storage"

	"charm.land/log/v2"
	"go.uber.org/fx"
)

func newServer(
	lc fx.Lifecycle,
	config *config.Config,
	catalogRepo *catalog.Repo,
	logger *log.Logger,
	contactHandler *contact.ContactUsHandler,
	ordersHandler *orders.Handler,
	userHandler *user.Handler,
) *server.HTTPServer {
	srv := server.NewHTTPServer(config, catalogRepo, contactHandler, ordersHandler, userHandler)

	lc.Append(fx.StartHook(func(ctx context.Context) {
		go func() {
			storage.Init(true) //TODO: move to own constructor

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
	)
}

func newContactHandler(emailService *email.Service, logger *log.Logger) *contact.ContactUsHandler {
	return contact.NewContactUsHandler(emailService, logger)
}

func newOrdersHandler(logger *log.Logger, emailService *email.Service, catalogRepo *catalog.Repo) *orders.Handler {
	return orders.NewHandler(logger, catalogRepo, emailService)
}

func newUserHandler() *user.Handler {
	return user.NewHandler()
}
