package storage

import (
	"database/sql"

	"github.com/charmbracelet/log"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {

	sqlConnectionString := viper.GetString("db-connection")
	var err error
	DB, err = gorm.Open(sqlite.Open(sqlConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", "error", err)
	}
	db, err := DB.DB()
	if err != nil {
		log.Fatal("failed to get database connection", "error", err)
	}
	err = migrateDB(db)
	if err != nil {
		log.Fatal("failed to migrate database", "error", err)
	}
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
