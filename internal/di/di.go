package di

import (
	"electrotech/internal/config"
	"electrotech/internal/email"
	"electrotech/internal/server"
	"fmt"
	"os"

	"charm.land/log/v2"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	app := fx.New(
		fx.Provide(newDB),
		fx.Provide(newCatalogRepo),
		fx.Provide(newOrdersRepo),
		fx.Provide(newUsersRepo),

		fx.Provide(newEmailService),
		fx.Provide(newConfig),
		fx.Provide(newServer),
		fx.Provide(newLogger),
		provideHandlers(),

		ftpModule(),

		fx.Invoke(func(_ *server.HTTPServer) {}),
	)

	return app
}

func newLogger() *log.Logger {
	logger := log.New(os.Stderr)
	logger.SetReportCaller(true)

	if os.Getenv("DEVEL") != "" {
		logger.SetLevel(log.DebugLevel)
	}

	return logger
}

func newConfig(logger *log.Logger) (*config.Config, error) {
	config, err := config.New(logger)
	if err != nil {
		return nil, fmt.Errorf("new config: %w", err)
	}

	return config, nil
}

func newEmailService(config *config.Config) *email.Service {
	return email.NewEmailService(config)
}
