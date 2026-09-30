package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"charm.land/log/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const StrongSecretMinLength = 20

var ErrSecretNotSet = errors.New("auth secret is not set")

func (h *Handler) getSecretKey() (string, error) {
	jwtSecret := h.config.Secret
	if jwtSecret == "" {
		h.logger.Error("auth secret isn't set")

		return "", ErrSecretNotSet
	}

	if len(jwtSecret) < StrongSecretMinLength {
		h.logger.Warn("auth secret is less than 20 characters, this must be security issue")
	}

	return jwtSecret, nil
}

const (
	TokenIssuer = "electrotech-back"
)

func (h *Handler) getKey() []byte {
	secret, err := h.getSecretKey()
	if err != nil {
		return nil
	}

	return []byte(secret)
}

type Claims struct {
	jwt.StandardClaims

	Email string `json:"email"`
	Id    int64  `json:"user_id"`
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})

			return
		}

		claims, err := h.ValidateToken(tokenString)
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

func (h *Handler) GenerateToken(email string, userID int64) (string, error) {
	expirationTime := time.Now().Add(h.config.TokenTTL)

	claims := &Claims{
		Email: email,
		Id:    userID,
		StandardClaims: jwt.StandardClaims{ //nolint:exhaustruct_v5
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expirationTime.Unix(),
			Issuer:    TokenIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(h.getKey())
	if err != nil {
		return "", fmt.Errorf("sign token string: %w", err)
	}

	return tokenString, nil
}

func (h *Handler) GenerateRefreshToken(userID int64) (string, error) {
	expirationTime := time.Now().Add(h.config.RefreshTokenTTL)

	claims := &Claims{ //nolint:exhaustruct_v5
		Id: userID,
		StandardClaims: jwt.StandardClaims{ //nolint:exhaustruct_v5
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expirationTime.Unix(),
			Issuer:    TokenIssuer,
			NotBefore: time.Now().Add(h.config.TokenTTL).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(h.getKey())
	if err != nil {
		return "", fmt.Errorf("sign token string: %w", err)
	}

	return tokenString, nil
}

func (h *Handler) ValidateToken(tokenString string) (*Claims, error) {
	var claims Claims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return h.getKey(), nil
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

	return &claims, nil
}
