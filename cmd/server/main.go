package main

import (
	"strings"
	"sync"

	"electrotech/internal/repository/catalog"
	ordersRepository "electrotech/internal/repository/orders"
	usersRepository "electrotech/internal/repository/users"
	"electrotech/internal/server"
	"electrotech/storage"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Warn("Can't load .env file", "error", err)
	}
	viper.SetConfigName("electrotech-back")
	viper.SetEnvKeyReplacer(
		strings.NewReplacer("-", "_", ".", "_"),
	)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("/etc/electrotech")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Warn("Failed read config file", "error", err)
	}
	log.SetReportCaller(true)

}

func main() {

	db, err := storage.New()
	if err != nil {
		log.Fatal("Can't init storage", "error", err)
	}

	usersRepo := usersRepository.New(db)
	catalogRepo, err := catalog.New()
	if err != nil {
		log.Fatalf("Error creating catalog repository: %v", err)
	}
	ordersRepo := ordersRepository.New(db)

	srv := server.NewHTTPServer(usersRepo, catalogRepo, ordersRepo)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() { mustRun(srv.Run(), &wg) }()
	if viper.GetBool("ftp.enable") {

		ftpServer, err := server.NewFTPServer()
		if err != nil {
			log.Fatal("Failed create FTP server", "error", err)
		}
		wg.Add(1)
		go func() { mustRun(ftpServer.Run(), &wg) }()
	}

	wg.Wait()
}

func mustRun(result error, wg *sync.WaitGroup) {
	defer wg.Done()
	if result != nil {
		log.Fatal("Failed run", "error", result)
	}
}
