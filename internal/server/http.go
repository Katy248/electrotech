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
)

const DefaultHTTPPort = 8080

type HTTPServer struct {
	engine *gin.Engine
	port   int
	logger *log.Logger
}

func NewHTTPServer(
	config *config.Config,
	logger *log.Logger,
	catalogRepo *catalog.Repo,
	contactHandler *contact.ContactUsHandler,
	ordersHandler *orders.Handler,
	userHandler *user.Handler,
	authHandler *auth.Handler,
	catalogHandler *catalogHandlers.Handler,
	catalogV2Handler *v2.Handler,
) *HTTPServer {
	gin.SetMode(config.GinMode)

	server := gin.Default()
	server.Use(newCORS())

	api := server.Group("/api")
	{
		api.Static("/files", config.Catalog.DataDir)
		api.POST("/contact-us", contactHandler.HandleContactUs())

		api.GET("/v2/products", catalogV2Handler.HandleGetProducts())
		{
			products := api.Group("/products")

			products.GET("/all/:page", catalogHandler.HandleGetProducts())
			products.POST("/filter/:page", catalogHandler.HandleGetProducts())
			products.GET("/:id", catalogHandler.HandleGetProduct())
		}
		{
			authGroup := api.Group("/auth")
			authGroup.POST("/login", authHandler.LoginHandler())
			authGroup.POST("/register", authHandler.RegisterHandler())
			authGroup.POST("/refresh", authHandler.Refresh())
		}
		{
			ordersGroup := api.Group("/orders")
			ordersGroup.Use(authHandler.AuthMiddleware())
			ordersGroup.POST("/create", ordersHandler.HandleCreateOrder())
			ordersGroup.GET("/get", ordersHandler.HandleGetUserOrders())
		}
		{
			usersGroup := api.Group("/user")
			usersGroup.Use(authHandler.AuthMiddleware())
			usersGroup.POST("/change-password", userHandler.HandleChangePassword())
			usersGroup.POST("/change-email", userHandler.HandleChangeEmail())
			usersGroup.POST("/change-phone", userHandler.HandleChangePhoneNumber())
			usersGroup.POST("/update-data", userHandler.HandleUpdateUserData())
			usersGroup.POST("/get-data", userHandler.HandleGetData())
			usersGroup.POST("/update-company-data", userHandler.HandleUpdateCompanyData())
			usersGroup.POST("/get-company-data", userHandler.HandleGetCompanyData())
		}
	}

	return &HTTPServer{engine: server, port: config.Port, logger: logger}
}

func (s *HTTPServer) Run() error {
	host := fmt.Sprintf(":%d", s.getPort())
	s.logger.Info("Starting server", "host", host)

	err := s.engine.Run(host)
	if err != nil {
		s.logger.Error("Failed run server", "error", err)

		return fmt.Errorf("run HTTP server: %w", err)
	}

	return nil
}

func (s *HTTPServer) getPort() int {
	var port = s.port
	if port <= 0 {
		s.logger.Warn("port value is invalid, fallback to default", "default", DefaultHTTPPort)
		port = DefaultHTTPPort
	}

	return port
}

func newCORS() gin.HandlerFunc {
	//nolint:exhaustruct_v5
	corsConf := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"POST", "GET", "OPTION", "DELETE", "PUT"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Origin", "X-Requested-With"},
		AllowCredentials: true,
	}

	return cors.New(corsConf)
}
