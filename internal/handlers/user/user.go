package user

import (
	"electrotech"
	"electrotech/internal/repository/users"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	gr "github.com/katy248/gravatar/pkg/url"
	"golang.org/x/crypto/bcrypt"
)

func ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ChangePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := users.ByEmail(c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword))
		if err != nil {
			log.Printf("Error comparing password: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		err = user.UpdatePassword(req.NewPassword)

		if err != nil {
			log.Error("Error hashing password", "error", err, "password", req.NewPassword)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		err = users.Update(user)
		if err != nil {
			log.Printf("Error updating password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
	}
}

func ChangeEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ChangeEmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := users.ByEmail(c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		user.Email = req.Email
		err = users.Update(user)
		if err != nil {
			log.Error("Error updating email", "error", err, "new-email", req.Email)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update email"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "email changed successfully"})
	}
}

func ChangePhoneNumber() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ChangePhoneNumberRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := users.ByEmail(c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Error("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		phone, err := electrotech.FormatPhoneNumber(req.PhoneNumber)
		if err != nil {
			log.Error("Error formatting phone number: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone number"})
			return
		}

		user.PhoneNumber = phone
		err = users.Update(user)

		if err != nil {
			log.Printf("Error updating phone number: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update phone number"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "phone number changed successfully"})
	}
}

func UpdateUserData() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateUserDataRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := users.ByEmail(c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		user.FirstName = req.FirstName
		user.LastName = req.LastName
		user.Surname = req.Surname
		err = users.Update(user)
		if err != nil {
			log.Printf("Error updating user data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User data updated successfully"})
	}
}

func GetData() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := users.ByEmail(c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Error("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
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
