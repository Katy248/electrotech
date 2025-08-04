package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	migration "electrotech/sql"

	_ "modernc.org/sqlite"
)

var (
	onlyLast = flag.Bool("last", false, "Only execute the last migration")
	index    = flag.Int("index", -1, "Execute migration at index")
)

func main() {
	flag.Parse()
	godotenv.Load()

	sqlConnectionString := os.Getenv("DB_CONNECTION")
	if sqlConnectionString == "" {
		log.Fatalf("DB_CONNECTION environment variable not set")
	}
	fmt.Println("Starting migration to " + sqlConnectionString)

	db, err := sql.Open("sqlite", sqlConnectionString)
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
