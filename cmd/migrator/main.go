package main

import (
	_ "embed"
	"flag"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"

	migration "electrotech/sql"
	"electrotech/storage"

	_ "modernc.org/sqlite"
)

var (
	onlyLast = flag.Bool("last", false, "Only execute the last migration")
	index    = flag.Int("index", -1, "Execute migration at index")
)

func main() {
	flag.Parse()
	if err := godotenv.Load(); err != nil {
		log.Error("Failed load .env", "error", err)
	}

	db, err := storage.New()
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if *onlyLast {
		if err := migration.MigrateLast(db); err != nil {
			log.Fatalf("Error executing migration: %v", err)
		}
	} else if *index >= 0 {
		if err := migration.MigrateIndex(db, *index); err != nil {
			log.Fatalf("Error executing migration: %v", err)
		}
	} else {
		if err := migration.Migrate(db); err != nil {
			log.Fatalf("Error executing migration: %v", err)
		}
	}

	fmt.Println("Migration executed successfully")

}
