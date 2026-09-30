package server

import (
	"electrotech/internal/config"
	"electrotech/internal/handlers/auth"
	catalogHandlers "electrotech/internal/handlers/catalog"
	v2 "electrotech/internal/handlers/catalog/v2"
	"electrotech/internal/handlers/contact"
	"electrotech/internal/handlers/orders"
	"electrotech/internal/handlers/user"
	"electrotech/internal/repository/catalog"
	"fmt"

	"charm.land/log/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HTTPServer struct {
	engine *gin.Engine
}

func NewHTTPServer(
	config *config.Config,
	catalogRepo *catalog.Repo,
	contactHandler *contact.ContactUsHandler,
	ordersHandler *orders.Handler,
	userHandler *user.Handler,
) *HTTPServer {
	gin.SetMode(config.GinMode)

	server := gin.Default()
	corsConf := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"POST", "GET", "OPTION", "DELETE", "PUT"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Origin", "X-Requested-With"},
		AllowCredentials: true,
	}
	server.Use(cors.New(corsConf))

	server.Use(func(ctx *gin.Context) {
		log.Info(ctx.Request.Header)
		ctx.Next()
	})

	api := server.Group("/api")
	{
		api.Static("/files", config.Catalog.DataDir)
		api.POST("/contact-us", contactHandler.HandleContactUs())

		api.GET("/v2/products", v2.GetProducts(catalogRepo))
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
			ordersGroup.POST("/create", ordersHandler.HandleCreateOrder())
			ordersGroup.GET("/get", ordersHandler.HandleGetUserOrders())
		}
		{
			usersGroup := api.Group("/user")
			usersGroup.Use(auth.AuthMiddleware())
			usersGroup.POST("/change-password", userHandler.HandleChangePassword())
			usersGroup.POST("/change-email", userHandler.HandleChangeEmail())
			usersGroup.POST("/change-phone", userHandler.HandleChangePhoneNumber())
			usersGroup.POST("/update-data", userHandler.HandleUpdateUserData())
			usersGroup.POST("/get-data", userHandler.HandleGetData())
			usersGroup.POST("/update-company-data", userHandler.HandleUpdateCompanyData())
			usersGroup.POST("/get-company-data", userHandler.HandleGetCompanyData())
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

		return fmt.Errorf("run HTTP server: %w", err)
	}

	return nil
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
