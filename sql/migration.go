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
)

var migrations = []string{migration1, migration2}

func Migrate(db *sql.DB) error {
	for i, migration := range migrations {
		log.Printf("Executing migration %d", i)
		fmt.Println(migration)
		_, err := db.Exec(migration)
		if err != nil {
			return err
		}
	}
	return nil
}
