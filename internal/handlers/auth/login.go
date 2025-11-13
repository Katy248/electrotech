package auth

import (
	"electrotech/internal/models"
	"electrotech/internal/users"
	"net/http"

	"github.com/charmbracelet/log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	Surname      string `json:"surname"`
	LastName     string `json:"last_name"`
	PhoneNumber  string `json:"phone_number"`
}

func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := users.ByEmail(req.Email)
		if err != nil || user.Email == "" {
			log.Printf("Error getting user by email '%s': %v", req.Email, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
		if err != nil {
			log.Printf("Error comparing password: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		response, err := getAuthResponse(user)
		if err != nil {
			log.Printf("Error generating auth response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		claimsUser, err := ValidateToken(req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		user, err := users.ByID(claimsUser.Id)
		if err != nil {
			log.Printf("Error getting user by id '%d': %v", claimsUser.Id, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		response, err := getAuthResponse(user)
		if err != nil {
			log.Printf("Error generating auth response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
func getAuthResponse(user *models.User) (*AuthResponse, error) {

	token, err := GenerateToken(user.Email, user.ID)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return nil, err
	}
	refreshToken, err := GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		return nil, err
	}

	return &AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		Email:        user.Email,
		FirstName:    user.FirstName,
		Surname:      user.Surname,
		LastName:     user.LastName,
		PhoneNumber:  user.PhoneNumber,
	}, nil
}
