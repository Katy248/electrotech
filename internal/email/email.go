package email

import (
	"fmt"
	"net/smtp"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Enabled  bool

	infoReceiver string
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
func (c *Config) Auth() smtp.Auth {
	return smtp.PlainAuth("", c.User, c.Password, c.Host)
}
func (c *Config) InfoReceiver() string {
	if c.infoReceiver == "" {
		return c.User
	}
	return c.infoReceiver
}

func getConfig() *Config {
	return &Config{
		Enabled:      viper.GetBool("mail.enable"),
		Port:         viper.GetInt("mail.port"),
		User:         viper.GetString("mail.user"),
		Password:     viper.GetString("mail.password"),
		Host:         viper.GetString("mail.host"),
		infoReceiver: viper.GetString("mail.info-receiver"),
	}
}

func IsEnabled() bool {
	conf := getConfig()

	return conf.Enabled
}

func Send(conf *Config, content []byte, to string) error {
	if !conf.Enabled {
		log.Error("Try to use mail system that not enabled")
		return fmt.Errorf("mail system not enabled")
	}

	err := smtp.SendMail(conf.Addr(), conf.Auth(), conf.User, []string{to}, content)
	if err != nil {
		log.Error("Failed to send email", "error", err)
		return err
	}
	return nil
}

func SendInfo(content []byte) error {
	conf := getConfig()
	return Send(conf, content, conf.InfoReceiver())
}
