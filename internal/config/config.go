package config

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Setup() {
	if err := godotenv.Load(); err != nil {
		log.Warn("Can't load .env file", "error", err)
	}
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
	log.SetReportCaller(true)

}
