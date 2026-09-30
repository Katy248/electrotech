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

type DBConfig struct {
	ConnectionString string `mapstructure:"connection-string"`
	AutoMigrate      bool   `mapstructure:"auto-migrate"`
}

func Connect(config DBConfig, logger *log.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.ConnectionString), &gorm.Config{}) //nolint:exhaustruct_v5
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if config.AutoMigrate {
		logger.Debug("Auto-migrating database")

		err := migrateDB(db)
		if err != nil {
			return nil, fmt.Errorf("migrate database: %w", err)
		}
	}

	return db, nil
}

func GetMigrationsDir() string {
	return "./sql/migrations"
}

func migrateDB(gormDB *gorm.DB) error {
	err := migration.Up(SQLConnection(gormDB), GetMigrationsDir())
	if err != nil {
		return fmt.Errorf("up migration: %w", err)
	}

	return nil
}
