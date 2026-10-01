package storage

import (
	"database/sql"
	"electrotech/storage/migration"
	"fmt"

	"charm.land/log/v2"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func SQLConnection(gormDB *gorm.DB) *sql.DB {
	db, err := gormDB.DB()
	if err != nil {
		log.Fatal("failed to get database connection", "error", err)
	}

	return db
}

type Config struct {
	ConnectionString string `mapstructure:"connection-string"`
	AutoMigrate      bool   `mapstructure:"auto-migrate"`
	MigrationsDir    string `mapstructure:"migrations-dir"`
}

func Connect(config Config, logger *log.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.ConnectionString), &gorm.Config{}) //nolint:exhaustruct_v5
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if config.AutoMigrate {
		logger.Info("Auto-migrating database")

		err := migrateDB(db, config.MigrationsDir)
		if err != nil {
			return nil, fmt.Errorf("migrate database: %w", err)
		}
	}

	return db, nil
}

func migrateDB(gormDB *gorm.DB, migrationsDir string) error {
	err := migration.Up(SQLConnection(gormDB), migrationsDir)
	if err != nil {
		return fmt.Errorf("up migration: %w", err)
	}

	return nil
}
