package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

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
	ftp "goftp.io/server/v2"
	"goftp.io/server/v2/driver/file"
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
		log.Fatal("Can't init storage", "error", err)
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
	ftpServer, err := newFTPServer()
	if err != nil {
		log.Fatal("Failed create FTP server", "error", err)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { mustRun(runServer(server), &wg) }()
	go func() { mustRun(runFTP(ftpServer), &wg) }()

	wg.Wait()
}

func mustRun(result error, wg *sync.WaitGroup) {
	defer wg.Done()
	if result != nil {
		log.Fatal("Failed run", "error", result)
	}
}

func newFTPServer() (*ftp.Server, error) {

	var conf struct {
		Enable   bool
		Port     int
		Username string
		Password string
		PublicIP string
	}

	conf.Port = viper.GetInt("ftp.port")
	conf.Enable = viper.GetBool("ftp.enable")
	conf.Username = viper.GetString("ftp.username")
	conf.Password = viper.GetString("ftp.password")
	conf.PublicIP = viper.GetString("ftp.public-ip")

	if !conf.Enable {
		return nil, nil
	}
	if conf.Port == 0 {
		return nil, fmt.Errorf("FTP port not specified")
	}
	if conf.Username == "" {
		return nil, fmt.Errorf("FTP user name not specified")
	}
	if conf.Password == "" {
		return nil, fmt.Errorf("FTP user password not specified")
	}
	if len(conf.Password) < 20 {
		log.Warn("FTP user password length is less than 20 symbols, this can be security issue")
	}
	ip := net.ParseIP(conf.PublicIP)
	if conf.PublicIP == "" {
		return nil, fmt.Errorf("FTP's public IP not specified")
	} else if ip.IsUnspecified() || ip.IsPrivate() {
		return nil, fmt.Errorf("bad FTP's public IP specification")
	}

	driver, err := file.NewDriver(viper.GetString("data-dir"))
	if err != nil {
		return nil, fmt.Errorf("failed create driver for data directory: %s", err)
	}

	usr := os.Getenv("USER")
	srv, err := ftp.NewServer(&ftp.Options{
		Driver: driver,
		Port:   conf.Port,
		Auth: &ftp.SimpleAuth{
			Name:     conf.Username,
			Password: conf.Password,
		},
		Perm:      ftp.NewSimplePerm(usr, usr),
		RateLimit: 1_000_000,
		PublicIP:  conf.PublicIP,
	})
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func runFTP(srv *ftp.Server) error {

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed serve FTP", "error", err)
		return err
	}
	return nil
}

func runServer(srv *gin.Engine) error {
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

	return err

}

func runSecure(srv *gin.Engine, host, cert, key string) error {
	srv.GET("/.well-known/acme-challenge/CYrYCTl0gu0IVWsQefp7m_CrUww7cNXf12p8IMz0sUk", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "CYrYCTl0gu0IVWsQefp7m_CrUww7cNXf12p8IMz0sUk.gp05os96QhdaYP8iPlcWow4JPYG8tW50-Pf3uzq5qiY")
	})
	return srv.RunTLS(host, cert, key)
}
