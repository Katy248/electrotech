package feedback

import (
	"bytes"
	"context"
	"electrotech/internal/models"
	_ "embed"
	"fmt"

	tmpl "html/template"

	"charm.land/log/v2"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type EmailService interface {
	SendInfo(content []byte, subject string) error
}

type Service struct {
	logger       *log.Logger
	emailService EmailService
	db           *gorm.DB
}

func New(
	logger *log.Logger,
	emailService EmailService,
	db *gorm.DB,

) *Service {
	return &Service{
		logger:       logger,
		emailService: emailService,
		db:           db,
	}
}

func (s *Service) NewUserQuestion(ctx context.Context, q *models.CreateUserQuestion) error {
	err := validator.New().Struct(q)
	if err != nil {
		return fmt.Errorf("invalid creation data: %w", err)
	}

	dbRequest := models.NewUserQuestion(
		q.PersonName, q.Email, q.Phone, q.Message, q.ClientIP,
	)

	err = s.db.WithContext(ctx).Create(&dbRequest).Error
	if err != nil {
		return fmt.Errorf("save to database: %w", err)
	}

	go s.sendEmail(dbRequest)

	return nil
}

//go:embed email.html
var EmailTemplate string

func (h *Service) sendEmail(question *models.UserQuestion) {
	content, err := h.buildEmail(question)
	if err != nil {
		h.logger.Error("Failed build email", "error", err)

		return
	}

	err = h.emailService.SendInfo(content, "Новый вопрос")
	if err != nil {
		h.logger.Error("Failed send email", "error", err)
	}
}

func (h *Service) buildEmail(question *models.UserQuestion) ([]byte, error) {
	template, err := tmpl.New("new-request-mail").Parse(EmailTemplate)
	if err != nil {
		h.logger.Error("Failed create email template for new request", "error", err)

		return nil, fmt.Errorf("failed create template: %w", err)
	}

	buff := &bytes.Buffer{}

	err = template.Execute(buff, question)
	if err != nil {
		return nil, fmt.Errorf("failed execute template: %w", err)
	}

	return buff.Bytes(), nil
}
