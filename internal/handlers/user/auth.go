package user

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	secretKey = "your-secret-key" // В реальном приложении используйте безопасный ключ из конфига
)

type Claims struct {
	Email  string `json:"email"`
	UserID int64  `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(email string, userID int64) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Email:  email,
		UserID: userID,
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

	return claims, nil
}
