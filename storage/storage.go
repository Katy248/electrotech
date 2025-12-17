package storage

import (
	"database/sql"
	"electrotech/storage/migration"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SQLConnection() *sql.DB {
	db, err := DB.DB()
	if err != nil {
		log.Fatal("failed to get database connection", "error", err)
	}
	return db
}

func Init(automigrate bool) {
	sqlConnectionString := viper.GetString("db-connection")
	var err error
	DB, err = gorm.Open(sqlite.Open(sqlConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", "error", err)
	}
	if automigrate {
		log.Debug("Auto-migrating database")
		migrateDB()
	}
}
func GetMigrationsDir() string {
	viper.SetDefault("migrations-dir", "./sql/migrations")
	return viper.GetString("migrations-dir")
}

func migrateDB() {
	err := migration.Up(SQLConnection(), GetMigrationsDir())
	if err != nil {
		log.Fatal("failed to migrate database", "error", err)
	}
}
