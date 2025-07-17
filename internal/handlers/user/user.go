package user

import (
	"electrotech/internal/repository/users"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

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

func LoginHandler(repo *users.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Error binding request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repo.GetByEmail(c.Request.Context(), req.Email)
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

		token, err := GenerateToken(user.Email, user.ID)
		if err != nil {
			log.Printf("Error generating token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		response := AuthResponse{
			Token:     token,
			Email:     user.Email,
			FirstName: user.FirstName,
			Surname:   user.Surname,
		}

		c.JSON(http.StatusOK, response)
	}
}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
