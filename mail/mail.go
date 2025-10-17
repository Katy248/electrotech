package mail

import (
	"fmt"
	"net/smtp"

	"github.com/charmbracelet/log"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
)

type Config struct {
	SenderName  string
	SenderEmail string
	SMTPServer  string
	SMTPPort    int

	AuthUsername string
	AuthIdentity string
	AuthPassword string
}

func (c *Config) From() string {
	return fmt.Sprintf("%s <%s>", c.SenderName, c.SenderEmail)
}

func (c *Config) SMTPHost() string {
	return fmt.Sprintf("%s:%d", c.SMTPServer, c.SMTPPort)
}

func (c *Config) Valid() bool {
	v := c != nil &&
		c.SenderName != "" &&
		c.SenderEmail != "" &&
		c.SMTPServer != "" &&
		c.SMTPPort > 0 &&
		c.AuthIdentity != "" &&
		c.AuthUsername != "" &&
		c.AuthPassword != ""
	if !v {
		log.Warn("Email configuration isn't valid. Skipping this functionality")
	}
	return v
}

func getConfig() *Config {
	return &Config{
		SenderName:   viper.GetString("email.sender-name"),
		SenderEmail:  viper.GetString("email.sender-email"),
		SMTPServer:   viper.GetString("email.smtp-server"),
		SMTPPort:     viper.GetInt("email.smtp-port"),
		AuthUsername: viper.GetString("email.username"),
		AuthIdentity: viper.GetString("email.identity"),
		AuthPassword: viper.GetString("email.password"),
	}
}

type Message struct {
	To          string
	Subject     string
	HTMLContent string
}

func Send(m Message) error {
	conf := getConfig()
	if !conf.Valid() {
		return nil
	}

	e := email.NewEmail()
	e.From = conf.From()
	e.To = []string{m.To}
	e.Subject = m.Subject
	e.HTML = []byte(m.HTMLContent)
	err := e.Send(conf.SMTPHost(), smtp.PlainAuth(conf.AuthIdentity, conf.AuthUsername, conf.AuthPassword, conf.SMTPServer))
	if err != nil {
		log.Error("Failed to send email", "to", m.To, "config", conf, "error", err)
		return err
	}
	return fmt.Errorf("not implemented")
}
