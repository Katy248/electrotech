package config

import (
	"electrotech/internal/repository/catalog"
	"electrotech/internal/server/ftp"
	"electrotech/storage"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"

	"charm.land/log/v2"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Devel     bool             `mapstructure:"devel"`
	GinMode   string           `mapstructure:"gin-mode"`
	Port      int              `mapstructure:"port"`
	JWTSecret string           `mapstructure:"jwt-secret"`
	Email     EmailConfig      `mapstructure:"mail"`
	Catalog   catalog.Config   `mapstructure:"catalog"`
	Auth      AuthConfig       `mapstructure:"auth"`
	DB        storage.DBConfig `mapstructure:"db"`
	FTP       ftp.Config       `mapstructure:"ftp"`
}

type AuthConfig struct {
	Secret          string        `mapstructure:"secret"`
	TokenTTL        time.Duration `mapstructure:"token-ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh-token-ttl"`
}

type EmailConfig struct {
	Host             string `mapstructure:"host"`
	Port             int    `mapstructure:"port"`
	User             string `mapstructure:"user"`
	Password         string `mapstructure:"password"`
	Enabled          bool   `mapstructure:"enabled"`
	InfoSender       string `mapstructure:"info-sender"`
	InfoReceiverConf string `mapstructure:"info-receiver"` // Conf postfix is fix for naming collision
}

func (c *EmailConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
func (c *EmailConfig) Auth() smtp.Auth {
	return smtp.PlainAuth("", c.User, c.Password, c.Host)
}
func (c *EmailConfig) InfoReceiver() string {
	if c.InfoReceiverConf == "" {
		return c.User
	}

	return c.InfoReceiverConf
}

func (c *EmailConfig) From() string {
	name := c.InfoSender
	if name == "" {
		name = "Electrotech info"
	}

	return fmt.Sprintf("%s <%s>", name, c.User)
}

func New(logger *log.Logger) (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Warn("Can't load .env file", "error", err)
	}

	viper.RegisterAlias("devel", "development")

	if os.Getenv("DEVEL") != "" {
		logger.Warn("Development mode enabled", "configFile", "electrotech-back.devel")
		viper.Set("devel", true)
		viper.SetConfigName("electrotech-back.devel")
	} else {
		viper.SetConfigName("electrotech-back")
	}

	viper.SetEnvKeyReplacer(
		strings.NewReplacer("-", "_", ".", "_"),
	)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("/etc/electrotech")
	viper.AutomaticEnv()

	// Default configurations
	viper.SetDefault("data-dir", "/data")

	err = viper.ReadInConfig()
	if err != nil {
		logger.Warn("Failed read config file", "error", err)
	}

	log.SetReportCaller(true)

	var config Config

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &config, nil
}
