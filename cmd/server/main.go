package main

import (
	"strings"
	"sync"

	"electrotech/internal/email"
	"electrotech/internal/repository/catalog"
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
	storage.Init()

	if email.IsEnabled() {
		log.Info("Mail system enabled")
	}

	catalogRepo, err := catalog.New()
	if err != nil {
		log.Fatalf("Error creating catalog repository: %v", err)
	}

	srv := server.NewHTTPServer(catalogRepo)

	var wg sync.WaitGroup
	wg.Add(1)
	go mustRun(srv.Run, &wg)
	if viper.GetBool("ftp.enable") {

		ftpServer, err := server.NewFTPServer()
		if err != nil {
			log.Fatal("Failed create FTP server", "error", err)
		}
		wg.Add(1)
		go mustRun(ftpServer.Run, &wg)
	}

	wg.Wait()
}

func mustRun(fn func() error, wg *sync.WaitGroup) {
	defer wg.Done()
	err := fn()
	if err != nil {
		log.Fatal("Failed run", "error", err)
	}
}
