package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	migration "electrotech/sql"

	_ "modernc.org/sqlite"
)

func main() {
	godotenv.Load()
	sqlConnectionString := os.Getenv("DB_CONNECTION")
	if sqlConnectionString == "" {
		log.Fatalf("DB_CONNECTION environment variable not set")
	}
	fmt.Println("Starting migration to " + sqlConnectionString)

	fmt.Println(migration.Migration)

	db, err := sql.Open("sqlite", sqlConnectionString)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if _, err := db.Exec(migration.Migration); err != nil {
		log.Fatalf("Error executing migration: %v", err)
	}

	fmt.Println("Migration executed successfully")

}
