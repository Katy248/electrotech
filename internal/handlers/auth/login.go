package auth

import (
	"electrotech"
	"electrotech/internal/models"
	"net/http"

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

func (h *Handler) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.logger.Error("Error binding login request body", "error", err)
			c.JSON(http.StatusBadRequest, electrotech.Error(err))

			return
		}

		user, err := h.usersRepo.ByEmail(c.Request.Context(), req.Email)
		if err != nil || user.Email == "" {
			h.logger.Errorf("Error getting user by email '%s': %v", req.Email, err)
			c.JSON(http.StatusUnauthorized, electrotech.ErrorStr("invalid credentials"))

			return
		}

		if !user.CheckPassword(req.Password) {
			h.logger.Error("Passwords don't match", "userId", user.ID)
			c.JSON(http.StatusUnauthorized, electrotech.ErrorStr("invalid credentials"))

			return
		}

		response, err := h.getAuthResponse(user)
		if err != nil {
			h.logger.Errorf("Error generating auth response: %v", err)
			c.JSON(http.StatusInternalServerError, electrotech.ErrorStr("invalid credentials"))

			return
		}

		c.JSON(http.StatusOK, response)
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
			c.JSON(http.StatusBadRequest, electrotech.Error(err))

			return
		}

		claimsUser, err := h.ValidateToken(req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, electrotech.Error(err))

			return
		}

		user, err := h.usersRepo.ByID(c.Request.Context(), claimsUser.Id)
		if err != nil {
			h.logger.Errorf("Error getting user by id '%d': %v", claimsUser.Id, err)
			c.JSON(http.StatusUnauthorized, electrotech.Error(err))

			return
		}

		response, err := h.getAuthResponse(user)
		if err != nil {
			h.logger.Errorf("Error generating auth response: %v", err)
			c.JSON(http.StatusInternalServerError, electrotech.Error(err))

			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func (h *Handler) getAuthResponse(user *models.User) (*AuthResponse, error) {
	token, err := h.GenerateToken(user.Email, user.ID)
	if err != nil {
		h.logger.Errorf("Error generating token: %v", err)

		return nil, err
	}

	refreshToken, err := h.GenerateRefreshToken(user.ID)
	if err != nil {
		h.logger.Errorf("Error generating refresh token: %v", err)

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
