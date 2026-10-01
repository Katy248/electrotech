//nolint:ireturn
package di

import (
	"context"
	"electrotech/internal/config"
	"electrotech/internal/server/ftp"
	"fmt"
	"os"

	"charm.land/log/v2"
	"go.uber.org/fx"
)

func ftpModule() fx.Option {
	if os.Getenv("FTP_ENABLED") == "" {
		return fx.Module("ftp")
	}

	return fx.Module("ftp",
		fx.Provide(ftpConfig),
		fx.Provide(newFTPServer),

		fx.Invoke(func(_ *ftp.Server) {}),
	)
}

func ftpConfig(config *config.Config) *ftp.Config {
	return &config.FTP
}

func newFTPServer(lc fx.Lifecycle, logger *log.Logger, config *ftp.Config) (*ftp.Server, error) {
	ftpServer, err := ftp.NewFTPServer(config, logger)
	if err != nil {
		return nil, fmt.Errorf("new ftp server: %w", err)
	}

	lc.Append(fx.StartHook(func(ctx context.Context) {
		go func() {
			err := ftpServer.Run()
			if err != nil {
				logger.Error("Failed run FTP server", "error", err)
			}
		}()
	}))

	return ftpServer, nil
}
