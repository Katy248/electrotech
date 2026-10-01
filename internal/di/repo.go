package di

import (
	"electrotech/internal/config"
	"electrotech/internal/repository/catalog"
	"electrotech/internal/repository/orders"
	"electrotech/internal/repository/users"
	"electrotech/storage"
	"fmt"

	"charm.land/log/v2"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func newDB(lc fx.Lifecycle, conf *config.Config, logger *log.Logger) (*gorm.DB, error) {
	db, err := storage.Connect(conf.DB, logger)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	lc.Append(fx.StopHook(func() error {
		db, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}

		err = db.Close()
		if err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}

		return nil
	}))

	return db, nil
}

func newCatalogRepo(conf *config.Config) (*catalog.Repo, error) {
	repo, err := catalog.New(&conf.Catalog)
	if err != nil {
		return nil, fmt.Errorf("new catalog repository: %w", err)
	}

	return repo, nil
}
func newOrdersRepo(db *gorm.DB, logger *log.Logger) *orders.Repo {
	return orders.NewRepo(db, logger)
}

func newUsersRepo(db *gorm.DB, logger *log.Logger) *users.Repo {
	return users.NewRepo(db, logger)
}
