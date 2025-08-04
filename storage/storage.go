package storage

import (
	"database/sql"
	"fmt"
	"os"

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
	return db, nil
}
