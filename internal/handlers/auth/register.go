package auth

import (
	"electrotech/internal/repository/users"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FirstName   string `json:"first_name" binding:"required"`
	Surname     string `json:"surname" binding:"required"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

func RegisterHandler(repo *users.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Проверяем, существует ли пользователь с таким email
		existingUser, err := repo.GetByEmail(c.Request.Context(), req.Email)
		if err == nil && existingUser.Email != "" {
			c.JSON(http.StatusConflict, gin.H{"error": "user with this email already exists"})
			return
		}

		// Хешируем пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		// Создаем нового пользователя
		err = repo.InsertNew(c.Request.Context(), users.InsertNewParams{
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
			FirstName:    req.FirstName,
			Surname:      req.Surname,
			LastName:     req.LastName,
			PhoneNumber:  req.PhoneNumber,
		})

		log.Printf("Error creating user: %v", err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
	}
}
