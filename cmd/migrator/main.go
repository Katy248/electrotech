// Migrator executable
package main

import (
	"electrotech/internal/config"
	"electrotech/storage"
	"electrotech/storage/migration"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
	logger "github.com/charmbracelet/log"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/viper"
)

var (
	log *logger.Logger
)

func main() {
	config.Setup()

	flagset := flag.NewFlagSet("migrator", flag.ExitOnError)

	debug := flagset.Bool("d", false, "Enable debug mode")
	if err := flagset.Parse(os.Args); err != nil {
		log.Error("Parse flags failed", "error", err)
		flag.Usage()
		os.Exit(1)
	}

	command := flagset.Arg(1)

	log = logger.New(os.Stderr)
	if *debug {
		log.SetLevel(logger.DebugLevel)
		log.Debug("Debug enabled")
	}

	log.SetTimeFormat("")

	switch command {
	case "up":
		up()
	case "down":
		down(flagset)
	case "info":
		info()
	case "":
		log.Error("Command not specified")
		usage(flagset)
	default:
		usage(flagset)
	}
}
func up() {
	log.Info("Migration UP")
	prepareDatabase()

	err := migration.Up(storage.SQLConnection(), storage.GetMigrationsDir())
	if err != nil {
		log.Error("Migration UP failed", "error", err)
	}
}
func down(set *flag.FlagSet) {
	migrationsCount := 1

	if set.NArg() >= 3 {
		var err error

		migrationsCount, err = strconv.Atoi(set.Arg(2))
		if err != nil {
			log.Error("Bad argument for down command, must be integer value", "error", err)
			os.Exit(1)
		}
	}

	log.Info("Migration DOWN", "by", migrationsCount)
	prepareDatabase()

	err := migration.DownBy(storage.SQLConnection(), storage.GetMigrationsDir(), migrationsCount)
	if err != nil {
		log.Error("Migration DOWN failed", "error", err)
	}
}
func info() {
	prepareDatabase()

	if records, err := migration.GetInformation(storage.SQLConnection()); err != nil {
		log.Error("Migration INFO failed", "error", err)
	} else {
		for _, r := range records {
			printRecord(r)
		}
	}
}

func printRecord(record *migrate.MigrationRecord) {
	title := lipgloss.NewStyle().Bold(true).Render(record.Id)
	migratedAt := lipgloss.NewStyle().Render(record.AppliedAt.Format(time.ANSIC))

	fmt.Println(title, migratedAt)
}

func prepareDatabase() {
	connectionString := viper.GetString("db-connection")
	if connectionString == "" {
		log.Error("Database connection string is not set")
	}

	storage.Init(false)
}
func usage(set *flag.FlagSet) {
	fmt.Printf("Usage: %s [options] <command>\n", set.Name())
	fmt.Printf("Commands:\n")
	fmt.Printf("  up\n")
	fmt.Printf("  down\n")
	fmt.Printf("  info\n")
	fmt.Printf("Options:\n")
	fmt.Printf("  -h - Print this help message\n")

	set.Usage()
}
