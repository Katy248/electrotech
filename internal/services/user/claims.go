package user

import (
	"electrotech/internal/models"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	TokenIssuer           = "electrotech-back"
	StrongSecretMinLength = 20
)

func (s *Service) GenerateToken(email string, userID int64) (string, error) {
	expirationTime := time.Now().Add(s.config.TokenTTL)

	claims := &models.Claims{
		Email: email,
		Id:    userID,
		StandardClaims: jwt.StandardClaims{ //nolint:exhaustruct_v5
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expirationTime.Unix(),
			Issuer:    TokenIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.getKey())
	if err != nil {
		return "", fmt.Errorf("sign token string: %w", err)
	}

	return tokenString, nil
}
func (s *Service) GenerateRefreshToken(userID int64) (string, error) {
	expirationTime := time.Now().Add(s.config.RefreshTokenTTL)

	claims := &models.Claims{ //nolint:exhaustruct_v5
		Id: userID,
		StandardClaims: jwt.StandardClaims{ //nolint:exhaustruct_v5
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expirationTime.Unix(),
			Issuer:    TokenIssuer,
			NotBefore: time.Now().Add(s.config.TokenTTL).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.getKey())
	if err != nil {
		return "", fmt.Errorf("sign token string: %w", err)
	}

	return tokenString, nil
}

var ErrSecretNotSet = errors.New("auth secret is not set")

func (s *Service) getSecretKey() (string, error) {
	jwtSecret := s.config.Secret
	if jwtSecret == "" {
		s.logger.Error("auth secret isn't set")

		return "", ErrSecretNotSet
	}

	if len(jwtSecret) < StrongSecretMinLength {
		s.logger.Warn("auth secret is less than 20 characters, this must be security issue")
	}

	return jwtSecret, nil
}

func (s *Service) getKey() []byte {
	secret, err := s.getSecretKey()
	if err != nil {
		return nil
	}

	return []byte(secret)
}

func (s *Service) newAuthResponse(user *models.User) (*models.AuthResponse, error) {
	token, err := s.GenerateToken(user.Email, user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	refreshToken, err := s.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &models.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		Email:        user.Email,
		FirstName:    user.FirstName,
		Surname:      user.Surname,
		LastName:     user.LastName,
		PhoneNumber:  user.PhoneNumber,
	}, nil
}

func (s *Service) ValidateToken(tokenString string) (*models.Claims, error) {
	var claims models.Claims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return s.getKey(), nil
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
