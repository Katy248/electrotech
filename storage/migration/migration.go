package migration

import (
	"database/sql"
	"fmt"

	logger "github.com/charmbracelet/log"
	migrate "github.com/rubenv/sql-migrate"
)

var log = logger.Default().WithPrefix("migration")

func Up(db *sql.DB, migrationsDir string) error {
	source := migrate.FileMigrationSource{
		Dir: migrationsDir,
	}
	log.Debug("Migrating database...")

	count, err := migrate.Exec(db, "sqlite3", source, migrate.Up)
	if err != nil {
		return err
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
		return err
	}
	log.Info("Database migrated", "migrations", count)
	return nil
}
func GetInformation(db *sql.DB) ([]*migrate.MigrationRecord, error) {
	records, err := migrate.GetMigrationRecords(db, "sqlite3")
	if err != nil {
		return nil, fmt.Errorf("failed get migration records: %s", err)
	}
	return records, nil
}
