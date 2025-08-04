package user

import (
	"electrotech/internal/handlers/auth"
	"electrotech/internal/repository/users"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ChangePassword(repo *users.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ChangePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repo.GetByEmail(c.Request.Context(), c.GetString("email"))
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

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		err = repo.UpdatePassword(c.Request.Context(), users.UpdatePasswordParams{
			Email:        c.GetString("email"),
			PasswordHash: string(hashedPassword),
		})
		if err != nil {
			log.Printf("Error updating password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
	}
}

func ChangeEmail(repo *users.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ChangeEmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repo.GetByEmail(c.Request.Context(), c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		err = repo.UpdateEmail(c.Request.Context(), users.UpdateEmailParams{
			Email: req.Email,
			ID:    user.ID,
		})
		if err != nil {
			log.Printf("Error updating email: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update email"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "email changed successfully"})

		log.Printf("// TODO: should implement refresh tokens")
		log.Printf("// TODO: should implement refresh tokens")
		log.Printf("// TODO: should implement refresh tokens")
		log.Printf("// TODO: should implement refresh tokens")
		log.Printf("// TODO: should implement refresh tokens")
	}
}

func ChangePhoneNumber(repo *users.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ChangePhoneNumberRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repo.GetByEmail(c.Request.Context(), c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		phone, err := auth.FormatPhoneNumber(req.PhoneNumber)
		if err != nil {
			log.Printf("Error formatting phone number: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone number"})
			return
		}

		err = repo.UpdatePhoneNumber(c.Request.Context(), users.UpdatePhoneNumberParams{
			PhoneNumber: phone,
			Email:       c.GetString("email"),
		})
		if err != nil {
			log.Printf("Error updating phone number: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update phone number"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "phone number changed successfully"})
	}
}

func UpdateUserData(repo *users.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateUserDataRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repo.GetByEmail(c.Request.Context(), c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		err = repo.UpdateData(c.Request.Context(), users.UpdateDataParams{
			Email:     c.GetString("email"),
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Surname:   req.Surname,
		})
		if err != nil {
			log.Printf("Error updating user data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User data updated successfully"})
	}
}

func GetData(repo *users.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := repo.GetByEmail(c.Request.Context(), c.GetString("email"))
		if err != nil || user.Email == "" {
			log.Printf("Error getting user by email '%s': %v", c.GetString("email"), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"first_name":   user.FirstName,
			"surname":      user.Surname,
			"last_name":    user.LastName,
		})
	}
}
