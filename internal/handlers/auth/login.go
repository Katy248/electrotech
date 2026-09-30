package auth

import (
	"electrotech"
	"electrotech/internal/models"
	"electrotech/internal/repository/users"
	"net/http"

	"charm.land/log/v2"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `binding:"required,email" json:"email"`
	Password string `binding:"required"       json:"password"`
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
			log.Error("Error binding login request body", "error", err)
			c.JSON(http.StatusBadRequest, electrotech.Error(err))

			return
		}

		user, err := users.ByEmail(req.Email)
		if err != nil || user.Email == "" {
			log.Errorf("Error getting user by email '%s': %v", req.Email, err)
			c.JSON(http.StatusUnauthorized, electrotech.ErrorStr("invalid credentials"))

			return
		}

		if !user.CheckPassword(req.Password) {
			log.Error("Passwords don't match", "userId", user.ID)
			c.JSON(http.StatusUnauthorized, electrotech.ErrorStr("invalid credentials"))

			return
		}

		response, err := getAuthResponse(user)
		if err != nil {
			log.Errorf("Error generating auth response: %v", err)
			c.JSON(http.StatusInternalServerError, electrotech.ErrorStr("invalid credentials"))

			return
		}

		c.JSON(http.StatusOK, response)
	}
}

type RefreshRequest struct {
	RefreshToken string `binding:"required" json:"refresh_token"`
}

func Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Errorf("Error binding request: %v", err)
			c.JSON(http.StatusBadRequest, electrotech.Error(err))

			return
		}

		claimsUser, err := ValidateToken(req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, electrotech.Error(err))

			return
		}

		user, err := users.ByID(claimsUser.Id)
		if err != nil {
			log.Errorf("Error getting user by id '%d': %v", claimsUser.Id, err)
			c.JSON(http.StatusUnauthorized, electrotech.Error(err))

			return
		}

		response, err := getAuthResponse(user)
		if err != nil {
			log.Errorf("Error generating auth response: %v", err)
			c.JSON(http.StatusInternalServerError, electrotech.Error(err))

			return
		}

		c.JSON(http.StatusOK, response)
	}
}
func getAuthResponse(user *models.User) (*AuthResponse, error) {
	token, err := GenerateToken(user.Email, user.ID)
	if err != nil {
		log.Errorf("Error generating token: %v", err)

		return nil, err
	}

	refreshToken, err := GenerateRefreshToken(user.ID)
	if err != nil {
		log.Errorf("Error generating refresh token: %v", err)

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
