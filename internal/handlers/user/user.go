package user

import (
	"electrotech"
	"electrotech/internal/repository/users"
	"net/http"

	"charm.land/log/v2"
	"github.com/gin-gonic/gin"
	gr "github.com/katy248/gravatar"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) HandleChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ChangePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding request: %v", err)
			c.JSON(http.StatusBadRequest, electrotech.Error(err))

			return
		}

		user, err := users.ByEmail(c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, electrotech.ErrorStr("invalid credentials"))

			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword))
		if err != nil {
			log.Printf("Error comparing password: %v", err)
			c.JSON(http.StatusUnauthorized, electrotech.Error(err))

			return
		}

		err = user.SetPassword(req.NewPassword)
		if err != nil {
			log.Error("Error hashing password", "error", err, "password", req.NewPassword)
			c.JSON(http.StatusInternalServerError, electrotech.ErrorStr("failed to hash password"))

			return
		}

		err = users.Update(user)
		if err != nil {
			log.Printf("Error updating password: %v", err)
			c.JSON(http.StatusInternalServerError, electrotech.ErrorStr("failed to update password"))

			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *Handler) HandleChangeEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ChangeEmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding request: %v", err)
			c.JSON(http.StatusBadRequest, electrotech.Error(err))

			return
		}

		user, err := users.ByEmail(c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, electrotech.ErrorStr("invalid credentials"))

			return
		}

		user.Email = req.Email

		err = users.Update(user)
		if err != nil {
			log.Error("Error updating email", "error", err, "new-email", req.Email)
			c.JSON(http.StatusInternalServerError, electrotech.ErrorStr("failed to update email"))

			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *Handler) HandleChangePhoneNumber() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ChangePhoneNumberRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding request: %v", err)
			c.JSON(http.StatusBadRequest, electrotech.Error(err))

			return
		}

		user, err := users.ByEmail(c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Error("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, electrotech.ErrorStr("invalid credentials"))

			return
		}

		phone, err := electrotech.FormatPhoneNumber(req.PhoneNumber)
		if err != nil {
			log.Error("Error formatting phone number: %v", err)
			c.JSON(http.StatusBadRequest, electrotech.Error(err))

			return
		}

		user.PhoneNumber = phone

		err = users.Update(user)
		if err != nil {
			log.Printf("Error updating phone number: %v", err)
			c.JSON(http.StatusInternalServerError, electrotech.ErrorStr("failed to update phone number"))

			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *Handler) HandleUpdateUserData() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateUserDataRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, electrotech.Error(err))

			return
		}

		user, err := users.ByEmail(c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, electrotech.ErrorStr("invalid credentials"))

			return
		}

		user.FirstName = req.FirstName
		user.LastName = req.LastName
		user.Surname = req.Surname

		err = users.Update(user)
		if err != nil {
			log.Printf("Error updating user data: %v", err)
			c.JSON(http.StatusInternalServerError, electrotech.ErrorStr("failed to update user data"))

			return
		}

		c.Status(http.StatusOK)
	}
}

func (h *Handler) HandleGetData() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := users.ByEmail(c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Error("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, electrotech.ErrorStr("invalid credentials"))

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"email":        user.Email,
			"avatarUrl":    gr.NewAvatarUrl(user.Email, gr.DefaultImage(gr.DefaultWavater)),
			"phone_number": user.PhoneNumber,
			"first_name":   user.FirstName,
			"surname":      user.Surname,
			"last_name":    user.LastName,
		})
	}
}
