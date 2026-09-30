package email

import (
	"electrotech/internal/config"
	"errors"
	"fmt"

	"charm.land/log/v2"
	e "github.com/jordan-wright/email"
)

type Service struct {
	Config *config.EmailConfig
}

func NewEmailService(config *config.Config) *Service {
	return &Service{
		Config: &config.Email,
	}
}

var ErrMailSystemNotEnabled = errors.New("mail system not enabled")

func (s *Service) SendInfo(content []byte, subject string) error {
	if !s.Config.Enabled {
		return ErrMailSystemNotEnabled
	}

	mail := e.NewEmail()
	mail.From = s.Config.From()
	mail.To = []string{s.Config.InfoReceiver()}
	mail.Subject = subject
	mail.HTML = content

	err := mail.Send(
		s.Config.Addr(),
		s.Config.Auth(),
	)
	if err != nil {
		log.Error("Failed send info email", "error", err, "mail", mail)

		return fmt.Errorf("failed send email: %w", err)
	}

	return nil
}
