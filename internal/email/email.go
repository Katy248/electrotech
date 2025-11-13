package email

import (
	"fmt"
	"net/smtp"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type Config struct {
	Host         string
	Port         int
	User         string
	Password     string
	Enabled      bool
	InfoSender   string
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

func (c *Config) From() string {
	name := c.InfoSender
	if name == "" {
		name = "Electrotech info"
	}

	return fmt.Sprintf("%s <%s>", name, c.Addr())
}

func getConfig() *Config {
	return &Config{
		Enabled:      viper.GetBool("mail.enable"),
		Port:         viper.GetInt("mail.port"),
		User:         viper.GetString("mail.user"),
		Password:     viper.GetString("mail.password"),
		Host:         viper.GetString("mail.host"),
		infoReceiver: viper.GetString("mail.info-receiver"),
		InfoSender:   viper.GetString("mail.info-sender"),
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

func SendInfo(content []byte, subject string) error {
	conf := getConfig()
	body := append(buildHeaders(conf, subject), content...)
	return Send(conf, body, conf.InfoReceiver())
}

func buildHeaders(conf *Config, subject string) []byte {
	return []byte(fmt.Sprintf(`
	Subject: %s
	From: %s
	Content-Type: text/html; charset="UTF-8"
	MIME-Version: 1.0
	Message-ID: <%s>


	`, subject, conf.From(), uuid.New().String()))
}
