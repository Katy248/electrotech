package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"electrotech/internal/handlers/auth"
	catalogHandlers "electrotech/internal/handlers/catalog"
	"electrotech/internal/handlers/filter"
	"electrotech/internal/handlers/orders"
	"electrotech/internal/handlers/user"
	"electrotech/internal/repository/catalog"
	ordersRepository "electrotech/internal/repository/orders"
	usersRepository "electrotech/internal/repository/users"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func init() { godotenv.Load() }

func getPort() int {
	var portStr = os.Getenv("PORT")
	if portStr == "" {
		log.Fatalf("PORT environment variable not set")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error parsing PORT environment variable: %v", err)
	}
	return port
}

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
	catalogRepo, err := catalog.New()
	ordersRepo := ordersRepository.New(db)

	server := gin.Default()
	// Enables cors
	corsConf := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"POST", "GET", "OPTION", "DELETE", "PUT"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Origin"},
		AllowCredentials: true,
	}
	server.Use(cors.New(corsConf))

	{
		api := server.Group("/api")
		{
			api.Static("/files", os.Getenv("DATA_DIR"))

			api.GET("/filters", filter.GetFilters(catalogRepo))

			if err != nil {
				log.Fatalf("Error creating catalog repository: %v", err)
			}

			{
				products := api.Group("/products")
				products.GET("/all/:page", catalogHandlers.GetProducts(catalogRepo))
				products.POST("/filter", catalogHandlers.GetProductsFilter(catalogRepo))
				products.GET("/:id", catalogHandlers.GetProduct(catalogRepo))
			}

			{
				authGroup := api.Group("/auth")
				authGroup.POST("/login", auth.LoginHandler(usersRepo))
				authGroup.POST("/register", auth.RegisterHandler(usersRepo))
				authGroup.POST("/refresh", auth.Refresh(usersRepo))
			}

			{
				ordersGroup := api.Group("/orders")
				ordersGroup.Use(auth.AuthMiddleware())

				ordersGroup.POST("/create", orders.CreateOrderHandler(ordersRepo, usersRepo, catalogRepo))
				ordersGroup.GET("/get", orders.GetUserOrdersHandler(ordersRepo))
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

	host := fmt.Sprintf(":%d", getPort())

	server.Run(host)
}
