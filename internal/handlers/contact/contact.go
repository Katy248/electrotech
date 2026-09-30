package contact

import (
	"bytes"
	"electrotech"
	"electrotech/internal/models"
	"electrotech/storage"
	"fmt"
	"net/http"
	"time"

	"charm.land/log/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	_ "embed"
	tmpl "html/template"
)

type EmailService interface {
	SendInfo(content []byte, subject string) error
}

type ContactUsHandler struct {
	emailService EmailService
	logger       *log.Logger
}

func NewContactUsHandler(emailService EmailService, logger *log.Logger) *ContactUsHandler {
	return &ContactUsHandler{
		emailService: emailService,
		logger:       logger,
	}
}

type Request struct {
	Name    string `binding:"required" json:"name"`
	Message string `binding:"required" json:"message"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

const DefautlRequestTimeout = time.Minute * 10

func GetRequestTimeout() time.Duration {
	viper.SetDefault("contact-us.request-timeout", DefautlRequestTimeout)

	return viper.GetDuration("contact-us.request-timeout")
}

func (h *ContactUsHandler) HandleContactUs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !h.checkRecentRequest(ip) {
			log.Warn("Too many requests", "clientIP", ip, "timeout", GetRequestTimeout())
			c.JSON(http.StatusTooManyRequests, electrotech.ErrorStr("too many requests"))

			return
		}

		var request Request
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, electrotech.Error(err))
			log.Error("Bad request", "error", err)

			return
		}

		if request.Email == "" && request.Phone == "" {
			c.JSON(http.StatusBadRequest, electrotech.ErrorStr("email or phone is required"))
			log.Error("Bad request", "error", "email or phone is required")

			return
		}

		dbRequest := models.NewUserQuestion(
			request.Name, request.Email, request.Phone, request.Message, ip,
		)

		err := storage.DB.Create(&dbRequest).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, electrotech.Error(err))
			log.Error("Failed create request", "error", err)
		}

		c.JSON(http.StatusOK, gin.H{})

		go h.sendEmail(dbRequest)
	}
}

//go:embed email.html
var EmailTemplate string

func (h *ContactUsHandler) sendEmail(question *models.UserQuestion) {
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

func (h *ContactUsHandler) buildEmail(question *models.UserQuestion) ([]byte, error) {
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

func (h *ContactUsHandler) checkRecentRequest(ip string) bool {
	var records []*models.UserQuestion

	dateAfter := time.Now().Add(-GetRequestTimeout())

	err := storage.DB.
		Model(new(models.UserQuestion)).
		Where("client_ip = ?", ip).
		Where("DATETIME(creation_date) > DATETIME(?)", dateAfter).
		Find(&records).Error
	if err != nil {
		h.logger.Error("Failed get recent requests", "error", err)

		return true
	}

	return len(records) == 0
}
