package sql

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
)

//go:embed migration-1.sql
var Migration string

//go:embed migration-2.sql
var Migration2 string

var migrations = []string{Migration, Migration2}

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
