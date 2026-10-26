package contact

import (
	"context"
	"electrotech"
	"electrotech/internal/models"
	"net/http"
	"time"

	"charm.land/log/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	_ "embed"
)

type ContactService interface {
	NewUserQuestion(ctx context.Context, q *models.CreateUserQuestion) error
}

type ContactUsHandler struct {
	logger  *log.Logger
	service ContactService
}

func NewContactUsHandler(
	logger *log.Logger,
	service ContactService,
) *ContactUsHandler {
	return &ContactUsHandler{
		logger:  logger,
		service: service,
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

		var request Request
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, electrotech.Error(err))
			log.Error("Bad request", "error", err)

			return
		}

		q := models.CreateUserQuestion{
			PersonName: request.Name,
			Email:      request.Email,
			Phone:      request.Phone,
			Message:    request.Message,
			ClientIP:   ip,
		}

		err := h.service.NewUserQuestion(c.Request.Context(), &q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, electrotech.Error(err))
			log.Error("Failed create request", "error", err)
		}

		c.JSON(http.StatusCreated, gin.H{})
	}
}
