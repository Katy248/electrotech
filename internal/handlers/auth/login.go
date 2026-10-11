package auth

import (
	"electrotech/internal/handlers/common"
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
			common.SendError(c, errBadRequest.WithDetails(err.Error()))

			return
		}

		user, err := h.usersRepo.ByEmail(c.Request.Context(), req.Email)
		if err != nil || user.Email == "" {
			h.logger.Errorf("Error getting user by email '%s': %v", req.Email, err)
			common.SendError(c, errInvalidCredentials)

			return
		}

		if !user.CheckPassword(req.Password) {
			h.logger.Error("Passwords don't match", "userId", user.ID)
			common.SendError(c, errInvalidCredentials)

			return
		}

		response, err := h.createAuthResponse(user)
		if err != nil {
			h.logger.Errorf("Error generating auth response: %v", err)
			common.SendError(c, err)

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
			common.SendError(c, errBadRequest.WithDetails(err.Error()))

			return
		}

		claimsUser, err := h.ValidateToken(req.RefreshToken)
		if err != nil {
			common.SendError(c, errInvalidToken.WithDetails(err.Error()))

			return
		}

		user, err := h.usersRepo.ByID(c.Request.Context(), claimsUser.Id)
		if err != nil {
			h.logger.Errorf("Error getting user by id '%d': %v", claimsUser.Id, err)
			common.SendError(c, errInvalidToken.WithDetails("no such user"))

			return
		}

		response, err := h.createAuthResponse(user)
		if err != nil {
			h.logger.Errorf("Error generating auth response: %v", err)
			common.SendError(c, err)

			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func (h *Handler) createAuthResponse(user *models.User) (*AuthResponse, error) {
	token, err := h.GenerateToken(user.Email, user.ID)
	if err != nil {
		return nil, errCreateAuthResponse.WithError(err)
	}

	refreshToken, err := h.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errCreateAuthResponse.WithError(err)
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
