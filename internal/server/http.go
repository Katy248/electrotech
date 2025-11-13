package server

import (
	"electrotech/internal/handlers/auth"
	catalogHandlers "electrotech/internal/handlers/catalog"
	"electrotech/internal/handlers/orders"
	"electrotech/internal/handlers/user"
	"electrotech/internal/repository/catalog"
	ordersRepo "electrotech/internal/repository/orders"
	"electrotech/internal/repository/users"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HTTPServer struct {
	engine *gin.Engine
}

func NewHTTPServer(usersRepo *users.Queries, catalogRepo *catalog.Repo, ordersRepo *ordersRepo.Queries) *HTTPServer {

	gin.SetMode(viper.GetString("gin-mode"))
	server := gin.Default()
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
			api.Static("/files", viper.GetString("data-dir"))

			{
				products := api.Group("/products")

				products.GET("/all/:page", catalogHandlers.GetProducts(catalogRepo))
				products.POST("/filter/:page", catalogHandlers.GetProducts(catalogRepo))
				products.GET("/:id", catalogHandlers.GetProduct(catalogRepo))
			}

			{
				authGroup := api.Group("/auth")
				authGroup.POST("/login", auth.LoginHandler())
				authGroup.POST("/register", auth.RegisterHandler())
				authGroup.POST("/refresh", auth.Refresh())
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
	return &HTTPServer{engine: server}
}

func (s *HTTPServer) Run() error {

	host := fmt.Sprintf(":%d", getPort())
	log.Info("Starting server", "host", host)

	err := s.engine.Run(host)

	if err != nil {
		log.Error("Failed run server", "error", err)
	}

	return err
}

const DefaultHTTPPort = 8080

func getPort() int {
	viper.SetDefault("port", DefaultHTTPPort)
	var port = viper.GetInt("port")
	if port == 0 {
		log.Warn("PORT value is invalid, fallback to default", "default", DefaultHTTPPort)
	}
	return port
}
