package auth

import (
	"context"
	"electrotech/internal/config"
	"electrotech/internal/handlers/common"
	"electrotech/internal/models"
	"net/http"

	"charm.land/log/v2"
	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	ByID(ctx context.Context, id int64) (*models.User, error)
	ByEmail(ctx context.Context, email string) (*models.User, error)
	InsertNew(ctx context.Context, user *models.User) error
}

type UserService interface {
	Login(ctx context.Context, email, password string) (*models.AuthResponse, error)
	Register(ctx context.Context, request *models.RegisterRequest) error
	Refresh(ctx context.Context, token string) (*models.AuthResponse, error)

	ValidateToken(tokenString string) (*models.Claims, error)
}

type Handler struct {
	logger    *log.Logger
	config    *config.AuthConfig
	usersRepo UserRepository
	service   UserService
}

func NewHandler(
	logger *log.Logger, config *config.AuthConfig, usersRepo UserRepository,
	service UserService,
) (*Handler, error) {
	handler := &Handler{
		logger:    logger,
		config:    config,
		usersRepo: usersRepo,
		service:   service,
	}

	return handler, nil
}

type LoginRequest struct {
	Email    string `binding:"required,email" json:"email"`
	Password string `binding:"required"       json:"password"`
}

func (h *Handler) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			common.SendError(c, errBadRequest.WithDetails(err.Error()))

			return
		}

		resp, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			common.SendError(c, err)

			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
func (h *Handler) RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			common.SendError(c, errBadRequest.WithDetails(err.Error()))

			return
		}

		err := h.service.Register(c.Request.Context(), &req)
		if err != nil {
			common.SendError(c, err)

			return
		}

		c.Status(http.StatusCreated)
	}
}

type RefreshRequest struct {
	RefreshToken string `binding:"required" json:"refresh_token"`
}

func (h *Handler) Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.logger.Errorf("Error binding request: %v", err)
			common.SendError(c, errBadRequest.WithDetails(err.Error()))

			return
		}

		response, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
		if err != nil {
			common.SendError(c, err)

			return
		}

		c.JSON(http.StatusOK, response)
	}
}
