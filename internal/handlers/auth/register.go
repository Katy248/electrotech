package auth

import (
	"electrotech"
	"electrotech/internal/models"
	"electrotech/internal/repository/users"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
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

func RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		req.Email = strings.ToLower(req.Email)

		// Проверяем, существует ли пользователь с таким email
		existingUser, err := users.ByEmail(req.Email)
		if err == nil && existingUser.Email != "" {
			log.Error("Attempt to create user with email already taken", "email", req.Email)
			c.JSON(http.StatusConflict, gin.H{"error": "user with this email already exists"})

			return
		}

		phone, err := electrotech.FormatPhoneNumber(req.PhoneNumber)
		if err != nil {
			log.Errorf("Error formatting phone number (is is probably invalid): %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone number"})

			return
		}

		user := &models.User{
			Email:       req.Email,
			FirstName:   req.FirstName,
			Surname:     req.Surname,
			LastName:    req.LastName,
			PhoneNumber: phone,
		}
		if err := user.SetPassword(req.Password); err != nil {
			log.Error("Failed set (hash) user password")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})

			return
		}

		// Создаем нового пользователя
		err = users.InsertNew(user)
		if err != nil {
			log.Errorf("Error creating user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})

			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
	}
}
