package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	TokenTTL        = time.Hour * 24
	RefreshTokenTTL = TokenTTL * 2
)

const (
	secretKey = "your-secret-key" // В реальном приложении используйте безопасный ключ из конфига
)

type Claims struct {
	Email string `json:"email"`
	Id    int64  `json:"user_id"`
	jwt.StandardClaims
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
		c.Set("user_id", claims.Id)
		c.Next()
	}
}

func GenerateToken(email string, userID int64) (string, error) {
	expirationTime := time.Now().Add(TokenTTL)

	claims := &Claims{
		Email: email,
		Id:    userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func GenerateRefreshToken(userID int64) (string, error) {
	expirationTime := time.Now().Add(RefreshTokenTTL)

	claims := &Claims{
		Id: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))

}
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	if err := claims.Valid(); err != nil {
		return nil, err
	}

	return claims, nil
}
