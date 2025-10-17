package storage

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	migrate "github.com/rubenv/sql-migrate"
	_ "modernc.org/sqlite"
)

func New() (*sql.DB, error) {
	sqlConnectionString := os.Getenv("DB_CONNECTION")

	if sqlConnectionString == "" {
		return nil, fmt.Errorf("DB_CONNECTION environment variable isn't set")
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
