package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	_secretKey = "" // В реальном приложении используйте безопасный ключ из конфига
)

const (
	TokenTTL        = time.Hour * 24
	RefreshTokenTTL = TokenTTL * 2
)

func Setup() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET environment variable isn't set")
	}
	if len(jwtSecret) < 20 {
		log.Warn("JWT_SECRET is less than 20 characters, this must be security issue")
	}
	_secretKey = jwtSecret
}

func getKey() []byte {
	return []byte(_secretKey)
}

type Claims struct {
	Email string `json:"email"`
	Id    int64  `json:"user_id"`
	jwt.StandardClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			log.Error("Failed to validate token", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
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
	return token.SignedString(getKey())
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
	return token.SignedString(getKey())

}
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return getKey(), nil
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
