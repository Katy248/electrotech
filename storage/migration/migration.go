package migration

import (
	"database/sql"
	"electrotech/sql/migrations"
	"fmt"

	"charm.land/log/v2"
	migrate "github.com/rubenv/sql-migrate"
)

func Up(db *sql.DB) error {
	source := migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrations.Migrations,
		Root:       ".",
	}

	log.Debug("Migrating database...")

	count, err := migrate.Exec(db, "sqlite3", source, migrate.Up)
	if err != nil {
		return fmt.Errorf("execute migration: %w", err)
	}

	log.Info("Database migrated", "migrations", count)

	return nil
}

func DownBy(db *sql.DB, migrationsDir string, migrationsCount int) error {
	source := migrate.FileMigrationSource{
		Dir: migrationsDir,
	}

	log.Debug("Migrating database...")

	count, err := migrate.ExecMax(db, "sqlite3", source, migrate.Down, migrationsCount)
	if err != nil {
		return fmt.Errorf("failed execute: %w", err)
	}

	log.Info("Database migrated", "migrations", count)

	return nil
}
func GetInformation(db *sql.DB) ([]*migrate.MigrationRecord, error) {
	records, err := migrate.GetMigrationRecords(db, "sqlite3")
	if err != nil {
		return nil, fmt.Errorf("failed get migration records: %w", err)
	}

	return records, nil
}
