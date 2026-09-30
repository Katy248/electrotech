package auth

import (
	"electrotech"
	"electrotech/internal/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email       string `binding:"required,email" json:"email"`
	Password    string `binding:"required,min=8" json:"password"`
	FirstName   string `binding:"required"       json:"first_name"`
	Surname     string `binding:"required"       json:"surname"`
	LastName    string `json:"last_name"`
	PhoneNumber string `binding:"required"       json:"phone_number"`
}

func (h *Handler) RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, electrotech.Error(err))

			return
		}

		req.Email = strings.ToLower(req.Email)

		// Проверяем, существует ли пользователь с таким email
		existingUser, err := h.usersRepo.ByEmail(req.Email)
		if err == nil && existingUser.Email != "" {
			h.logger.Error("Attempt to create user with email already taken", "email", req.Email)
			c.JSON(http.StatusConflict, electrotech.Error(err))

			return
		}

		phone, err := electrotech.FormatPhoneNumber(req.PhoneNumber)
		if err != nil {
			h.logger.Errorf("Error formatting phone number (is is probably invalid): %v", err)
			c.JSON(http.StatusBadRequest, electrotech.Error(err))

			return
		}

		user := &models.User{ //nolint:exhaustruct_v5
			Email:       req.Email,
			FirstName:   req.FirstName,
			Surname:     req.Surname,
			LastName:    req.LastName,
			PhoneNumber: phone,
		}
		if err := user.SetPassword(req.Password); err != nil {
			h.logger.Error("Failed set (hash) user password")
			c.JSON(http.StatusInternalServerError, electrotech.ErrorStr("failed set password"))

			return
		}

		// Создаем нового пользователя
		err = h.usersRepo.InsertNew(user)
		if err != nil {
			h.logger.Errorf("Error creating user: %v", err)
			c.JSON(http.StatusInternalServerError, electrotech.Error(err))

			return
		}

		c.Status(http.StatusCreated)
	}
}
