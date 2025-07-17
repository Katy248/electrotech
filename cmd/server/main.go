package main

import (
	"database/sql"
	"log"
	"os"

	catalogHandlers "electrotech/internal/handlers/catalog"
	"electrotech/internal/handlers/user"
	"electrotech/internal/repository/catalog"
	usersRepository "electrotech/internal/repository/users"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func init() { godotenv.Load() }

func main() {
	sqlConnectionString := os.Getenv("DB_CONNECTION")
	if sqlConnectionString == "" {
		log.Fatalf("DB_CONNECTION environment variable not set")
	}

	db, err := sql.Open("sqlite", sqlConnectionString)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	usersRepo := usersRepository.New(db)

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

			usersGroup := api.Group("/user")
			usersGroup.POST("/register", user.RegisterHandler(usersRepo))
			usersGroup.POST("/login", user.LoginHandler(usersRepo))
		}

	}

	server.Run(":1488")
}
