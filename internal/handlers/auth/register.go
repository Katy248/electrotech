package auth

import (
	"electrotech/internal/repository/users"
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
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

func FormatPhoneNumber(phone string) (string, error) {
	if phone == "" {
		return "", fmt.Errorf("phone number is empty")

	}
	if len(phone) < 11 {
		return "", fmt.Errorf("probably invalid phone number, length is less than 11 (%d)", len(phone))
	}
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	phone = strings.ReplaceAll(phone, " ", "")

	if phone[0] == '8' {
		phone = "+7" + phone[1:]
	}

	for index, ch := range phone {
		if ch != '+' && ch < '0' || ch > '9' {
			return phone, fmt.Errorf("invalid character %q at index %d", ch, index)
		}
	}

	if len(phone) != 12 {
		return phone, fmt.Errorf("invalid phone number length %d", len(phone))
	}
	return phone, nil
}

func GeneratePasswordHash(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(pass, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
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
		hashedPassword, err := GeneratePasswordHash(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		phone, err := FormatPhoneNumber(req.PhoneNumber)
		if err != nil {
			log.Printf("Error formatting phone number (is is probably invalid): %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone number"})
			return
		}

		// Создаем нового пользователя
		err = repo.InsertNew(c.Request.Context(), users.InsertNewParams{
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
			FirstName:    req.FirstName,
			Surname:      req.Surname,
			LastName:     req.LastName,
			PhoneNumber:  phone,
		})

		log.Printf("Error creating user: %v", err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
	}
}
