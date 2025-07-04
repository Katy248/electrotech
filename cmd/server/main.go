package main

import (
	"log"
	"os"

	catalogHandlers "electrotech/internal/handlers/catalog"
	"electrotech/internal/repository/catalog"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() { godotenv.Load() }

func main() {
	sqlConnectionString := os.Getenv("DB_CONNECTION")
	if sqlConnectionString == "" {
		log.Fatalf("DB_CONNECTION environment variable not set")
	}

	server := gin.Default()
	{
		api := server.Group("/api")
		{
			catalogRepo, err := catalog.New()
			if err != nil {
				log.Fatalf("Error creating catalog repository: %v", err)
			}
			products := api.Group("/products")
			products.GET("/all", catalogHandlers.GetProducts(catalogRepo))
		}

	}

	server.Run(":1488")
}
