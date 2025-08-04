package sql

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
)

var (
	//go:embed migrations/migration-1.sql
	migration1 string
	//go:embed migrations/migration-2.sql
	migration2 string
	//go:embed migrations/migration-3.sql
	migration3 string
)

var migrations = []string{migration1, migration2, migration3}

func runMigration(db *sql.DB, migration string) error {
	fmt.Println()
	fmt.Println(migration)
	_, err := db.Exec(migration)
	return err
}

func MigrateLast(db *sql.DB) error {
	return runMigration(db, migrations[len(migrations)-1])
}
func MigrateIndex(db *sql.DB, index int) error {
	if index >= len(migrations) {
		return fmt.Errorf("migration index out of range")
	}
	return runMigration(db, migrations[index])
}

func Migrate(db *sql.DB) error {
	for i, migration := range migrations {
		log.Printf("Executing migration %d", i)
		err := runMigration(db, migration)
		if err != nil {
			return err
		}
	}
	return nil
}
