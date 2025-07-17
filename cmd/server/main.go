package main

import (
	"database/sql"
	"log"
	"os"

	"electrotech/internal/handlers/auth"
	catalogHandlers "electrotech/internal/handlers/catalog"
	"electrotech/internal/handlers/user"
	"electrotech/internal/repository/catalog"
	usersRepository "electrotech/internal/repository/users"

	"github.com/gin-contrib/cors"
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
	// Enables cors
	server.Use(cors.Default())

	{
		api := server.Group("/api")
		{
			catalogRepo, err := catalog.New()
			if err != nil {
				log.Fatalf("Error creating catalog repository: %v", err)
			}

			{
				products := api.Group("/products")
				products.GET("/all", catalogHandlers.GetProducts(catalogRepo))
			}

			{
				authGroup := api.Group("/auth")
				authGroup.POST("/login", auth.LoginHandler(usersRepo))
				authGroup.POST("/register", auth.RegisterHandler(usersRepo))
				authGroup.POST("/refresh", auth.Refresh(usersRepo))
			}

			{
				usersGroup := api.Group("/user")
				usersGroup.Use(auth.AuthMiddleware())

				usersGroup.POST("/change-password", user.ChangePassword(usersRepo))
				usersGroup.POST("/change-email", user.ChangeEmail(usersRepo))
				usersGroup.POST("/change-phone", user.ChangePhoneNumber(usersRepo))
				usersGroup.POST("/update-data", user.UpdateUserData(usersRepo))
				usersGroup.POST("/get-data", user.GetData(usersRepo))
				usersGroup.POST("/update-company-data", user.UpdateCompanyData(usersRepo))
				usersGroup.POST("/get-company-data", user.GetCompanyData(usersRepo))
			}
		}

	}

	server.Run(":1488")
}
