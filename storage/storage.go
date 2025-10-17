package storage

import (
	"database/sql"
	"fmt"

	"github.com/charmbracelet/log"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
)

func New() (*sql.DB, error) {
	sqlConnectionString := viper.GetString("db-connection")

	if sqlConnectionString == "" {
		return nil, fmt.Errorf("db-connection value isn't set")
	}

	db, err := sql.Open("sqlite", sqlConnectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := migrateDB(db); err != nil {
		return nil, fmt.Errorf("failed migrate database: %s", err)
	}

	return db, nil
}

func migrateDB(db *sql.DB) error {
	source := migrate.FileMigrationSource{
		Dir: "./sql/migrations",
	}
	log.Debug("Migrating database...")
	count, err := migrate.Exec(db, "sqlite3", source, migrate.Up)
	if err != nil {
		return err
	}
	log.Info("Database migrated", "migrations", count)
	return nil
}
