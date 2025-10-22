package main

import (
	"fmt"
	"net/http"
	"strings"

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
	"github.com/spf13/viper"
)

func init() {
	godotenv.Load()
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
}

const DefaultPort = 8080

func getPort() int {
	var port = viper.GetInt("port")
	if port == 0 {
		log.Warn("PORT value is not set, fallback to default", "default", DefaultPort)
		return DefaultPort
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

	gin.SetMode(viper.GetString("gin-mode"))
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
			api.Static("/files", viper.GetString("data-dir"))

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

	runServer(server)
}

func runServer(srv *gin.Engine) {
	host := fmt.Sprintf(":%d", getPort())
	log.Info("Starting server", "host", host)

	tlsCert := viper.GetString("tls-cert")
	tlsKey := viper.GetString("tls-key")

	var err error

	if tlsCert != "" && tlsKey != "" {
		err = runSecure(srv, host, tlsCert, tlsKey)
	} else {
		log.Warn("Run insecure")
		err = srv.Run(host)
	}

	if err != nil {
		log.Error("Failed run server", "error", err)
	}

}

func runSecure(srv *gin.Engine, host, cert, key string) error {
	srv.GET("/.well-known/acme-challenge/CYrYCTl0gu0IVWsQefp7m_CrUww7cNXf12p8IMz0sUk", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "CYrYCTl0gu0IVWsQefp7m_CrUww7cNXf12p8IMz0sUk.gp05os96QhdaYP8iPlcWow4JPYG8tW50-Pf3uzq5qiY")
	})
	return srv.RunTLS(host, cert, key)
}
