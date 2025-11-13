package email

import (
	"fmt"
	"net/smtp"

	"github.com/charmbracelet/log"
	e "github.com/jordan-wright/email"
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

func SendInfo(content []byte, subject string) error {
	conf := getConfig()
	if !conf.Enabled {
		return fmt.Errorf("mail system not enabled")
	}
	mail := e.NewEmail()
	mail.From = conf.From()
	mail.To = []string{conf.InfoReceiver()}
	mail.Subject = subject
	mail.HTML = content

	err := mail.Send(
		conf.Addr(),
		conf.Auth(),
	)
	if err != nil {
		log.Error("Failed send info email", "error", err, "mail", mail)
		return fmt.Errorf("failed send email: %s", err)
	}
	return nil
}
