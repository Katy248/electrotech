package handlers

import (
	"context"
	"electrotech/internal/config"
	"electrotech/internal/gen/handler_gen"
	"electrotech/internal/models"
	"errors"
	"fmt"
)

type Validator interface {
	ValidateToken(token string) (*models.Claims, error)
}

type AuthMiddleware struct {
	config  *config.AuthConfig
	service Validator
}

var _ handler_gen.SecurityHandler = (*AuthMiddleware)(nil)

func NewAuthMiddleware(
	config *config.AuthConfig,
	service Validator,
) *AuthMiddleware {
	return &AuthMiddleware{
		config:  config,
		service: service,
	}
}

var (
	errAuth         = errors.New("auth error")
	errNoKey        = fmt.Errorf("%w: no auth token provided", errAuth)
	errInvalidToken = fmt.Errorf("%w: invalid token", errAuth)
)

type ctxKey string

const (
	userIDKey = ctxKey("user_id")
	emailKey  = ctxKey("email")
)

func GetUserID(ctx context.Context) (int64, bool) {
	val := ctx.Value(userIDKey)

	id, isInt := val.(int64)
	if !isInt {
		return 0, false
	}

	return id, true
}

// HandleJwtAuth implements [handler_gen.SecurityHandler].
func (a *AuthMiddleware) HandleJwtAuth(
	ctx context.Context,
	operationName handler_gen.OperationName,
	t handler_gen.JwtAuth,
) (context.Context, error) {
	token := t.GetAPIKey()
	if token == "" {
		return ctx, errNoKey
	}

	claims, err := a.service.ValidateToken(token)
	if err != nil {
		return ctx, errInvalidToken
	}

	ctx = context.WithValue(ctx, emailKey, claims.Email)
	ctx = context.WithValue(ctx, userIDKey, claims.Id)

	return ctx, nil
}
