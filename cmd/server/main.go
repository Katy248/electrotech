package main

import (
	"fmt"
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
	"electrotech/storage"

	"github.com/charmbracelet/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() { godotenv.Load() }

const defaultPort = 8080

func getPort() int {
	var portStr = os.Getenv("PORT")
	if portStr == "" {
		log.Warn("PORT environment variable not set, fallback to default", "default", defaultPort)
		return defaultPort
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Error("Failed parse PORT environment variable Fallback to default %d", "value", portStr, "error", err, "default", defaultPort)
	}
	return port
}

func main() {

	db, err := storage.New()
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	auth.Setup()

	usersRepo := usersRepository.New(db)
	catalogRepo, err := catalog.New()
	ordersRepo := ordersRepository.New(db)

	gin.SetMode(os.Getenv("GIN_MODE"))
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
				ordersGroup.GET("/get", orders.GetUserOrdersHandler(ordersRepo, catalogRepo))
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
