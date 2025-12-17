package contact

import (
	"bytes"
	"electrotech/internal/email"
	"electrotech/internal/models"
	"electrotech/storage"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	_ "embed"
	tmpl "html/template"
)

type Request struct {
	Name    string `json:"name" binding:"required"`
	Message string `json:"message" binding:"required"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

func GetRequestTimeout() time.Duration {
	viper.SetDefault("contact-us.request-timeout", time.Minute*10)
	return viper.GetDuration("contact-us.request-timeout")
}

func checkRecentRequest(ip string) bool {
	var records []*models.UserQuestion

	dateAfter := time.Now().Add(-GetRequestTimeout())

	if err := storage.DB.
		Model(&models.UserQuestion{}).
		Where("client_ip = ?", ip).
		Where("DATETIME(creation_date) > DATETIME(?)", dateAfter).
		Find(&records).Error; err != nil {
		log.Error("Failed get recent requests", "error", err)
		return true
	}
	return len(records) == 0
}

func ContactUsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		ip := c.ClientIP()
		if !checkRecentRequest(ip) {
			log.Warn("Too many requests", "clientIP", ip, "timeout", GetRequestTimeout())
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}

		var request Request
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Error("Bad request", "error", err)
			return
		}

		if request.Email == "" && request.Phone == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email or phone is required"})
			log.Error("Bad request", "error", "email or phone is required")
			return
		}

		dbRequest := models.UserQuestion{
			CreationDate: time.Now(),
			PersonName:   request.Name,
			Email:        &request.Email,
			Phone:        &request.Phone,
			Message:      request.Message,
			ClientIP:     ip,
		}

		err := storage.DB.Create(&dbRequest).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error("Failed create request", "error", err)
		}

		c.JSON(http.StatusOK, gin.H{})
		go sendEmail(&dbRequest)
	}
}

//go:embed email.html
var EmailTemplate string

func sendEmail(question *models.UserQuestion) {
	content, err := buildEmail(question)
	if err != nil {
		log.Error("Failed build email", "error", err)
		return
	}
	err = email.SendInfo(content, "Новый вопрос")
	if err != nil {
		log.Error("Failed send email", "error", err)
	}
}

func buildEmail(question *models.UserQuestion) ([]byte, error) {
	template, err := tmpl.New("new-request-mail").Parse(EmailTemplate)
	if err != nil {
		log.Error("Failed create email template for new request", "error", err)
		return nil, fmt.Errorf("failed create template: %s", err)
	}
	buff := &bytes.Buffer{}
	err = template.Execute(buff, question)
	if err != nil {
		return nil, fmt.Errorf("failed execute template: %s", err)
	}
	return buff.Bytes(), nil
}
