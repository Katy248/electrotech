package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const StrongSecretMinLength = 20

func getSecretKey() string {
	jwtSecret := viper.GetString("jwt-secret")
	if jwtSecret == "" {
		log.Fatal("jwt-secret isn't set")
	}

	if len(jwtSecret) < StrongSecretMinLength {
		log.Warn("jwt-secret is less than 20 characters, this must be security issue")
	}

	return jwtSecret
}

const (
	// TokenTTL - token time to live.
	//
	// TODO: make configurable.
	TokenTTL        = time.Hour * 2
	RefreshTokenTTL = time.Hour * 48 // TODO: make configurable
	TokenIssuer     = "electrotech-back"
)

func getKey() []byte {
	return []byte(getSecretKey())
}

type Claims struct {
	jwt.StandardClaims

	Email string `json:"email"`
	Id    int64  `json:"user_id"`
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
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expirationTime.Unix(),
			Issuer:    TokenIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(getKey())
	if err != nil {
		return "", fmt.Errorf("sign token string: %w", err)
	}

	return tokenString, nil
}

func GenerateRefreshToken(userID int64) (string, error) {
	expirationTime := time.Now().Add(RefreshTokenTTL)

	claims := &Claims{
		Id: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expirationTime.Unix(),
			Issuer:    TokenIssuer,
			NotBefore: time.Now().Add(TokenTTL).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(getKey())
	if err != nil {
		return "", fmt.Errorf("sign token string: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return getKey(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token invalid: %w", jwt.ErrSignatureInvalid)
	}

	if err := claims.Valid(); err != nil {
		return nil, fmt.Errorf("claims invalid: %w", err)
	}

	return claims, nil
}
