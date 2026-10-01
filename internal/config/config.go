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
	Devel   bool           `mapstructure:"devel"`
	GinMode string         `mapstructure:"gin-mode"`
	Port    int            `mapstructure:"port"`
	Email   EmailConfig    `mapstructure:"mail"`
	Catalog catalog.Config `mapstructure:"catalog"`
	Auth    AuthConfig     `mapstructure:"auth"`
	DB      storage.Config `mapstructure:"db"`
	FTP     ftp.Config     `mapstructure:"ftp"`
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

	viper.SetEnvPrefix("EL")
	viper.SetEnvKeyReplacer(
		strings.NewReplacer("-", "_", ".", "_"),
	)
	viper.AutomaticEnv()

	viper.AddConfigPath(".")
	viper.AddConfigPath("/app")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("/etc/electrotech")

	// Default configurations
	viper.SetDefault("data-dir", "/data")

	err = viper.ReadInConfig()
	if err != nil {
		logger.Warn("Failed read config file", "error", err)
	}

	setDefaults(viper.GetViper())

	config := unmarshalConfig(viper.GetViper())

	return config, nil
}

const (
	DefaultHTTPPort = 8080
	DefaultFTPPort  = 8021
)

func setDefaults(viper *viper.Viper) {
	viper.SetDefault("port", DefaultHTTPPort)
	viper.SetDefault("ftp.port", DefaultFTPPort)
}

func unmarshalConfig(viper *viper.Viper) *Config {
	return &Config{
		Devel:   viper.GetBool("devel"),
		GinMode: viper.GetString("gin-mode"),
		Port:    viper.GetInt("port"),
		Email: EmailConfig{
			Host:             viper.GetString("mail.host"),
			Port:             viper.GetInt("mail.port"),
			User:             viper.GetString("mail.user"),
			Password:         viper.GetString("mail.password"),
			Enabled:          viper.GetBool("mail.enabled"),
			InfoSender:       viper.GetString("mail.info-sender"),
			InfoReceiverConf: viper.GetString("mail.info-receiver"),
		},
		Catalog: catalog.Config{
			DataDir: viper.GetString("catalog.data-dir"),
		},
		Auth: AuthConfig{
			Secret:          viper.GetString("auth.secret"),
			TokenTTL:        viper.GetDuration("auth.token-ttl"),
			RefreshTokenTTL: viper.GetDuration("auth.refresh-token-ttl"),
		},
		DB: storage.Config{
			ConnectionString: viper.GetString("db.connection-string"),
			AutoMigrate:      viper.GetBool("db.auto-migrate"),
		},
		FTP: ftp.Config{
			Port:     viper.GetInt("ftp.port"),
			Username: viper.GetString("ftp.username"),
			Password: viper.GetString("ftp.password"),
			PublicIP: viper.GetString("ftp.public-ip"),
		},
	}
}
