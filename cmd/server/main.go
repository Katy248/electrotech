// Server executable
package main

import (
	"sync"

	"electrotech/internal/config"
	"electrotech/internal/email"
	"electrotech/internal/repository/catalog"
	"electrotech/internal/server"
	"electrotech/storage"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func main() {
	config.Setup()
	storage.Init(true)

	if email.IsEnabled() {
		log.Info("Mail system enabled")
	}

	catalogRepo, err := catalog.New()
	if err != nil {
		log.Fatalf("Error creating catalog repository: %v", err)
	}

	srv := server.NewHTTPServer(catalogRepo)

	var waitGroup sync.WaitGroup

	waitGroup.Go(func() {
		err := srv.Run()
		if err != nil {
			log.Error("Failed run", "error", err)
		}
	})

	if viper.GetBool("ftp.enable") {
		ftpServer, err := server.NewFTPServer()
		if err != nil {
			log.Fatal("Failed create FTP server", "error", err)
		}

		waitGroup.Go(func() {
			err := ftpServer.Run()
			if err != nil {
				log.Error("Failed run", "error", err)
			}
		})
	}

	waitGroup.Wait()
}
