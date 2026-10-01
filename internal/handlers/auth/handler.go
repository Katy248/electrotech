package auth

import (
	"context"
	"electrotech/internal/config"
	"electrotech/internal/models"
	"fmt"

	"charm.land/log/v2"
)

type UserRepository interface {
	ByID(ctx context.Context, id int64) (*models.User, error)
	ByEmail(ctx context.Context, email string) (*models.User, error)
	InsertNew(ctx context.Context, user *models.User) error
}

type Handler struct {
	logger    *log.Logger
	config    *config.AuthConfig
	usersRepo UserRepository
}

func NewHandler(logger *log.Logger, config *config.AuthConfig, usersRepo UserRepository) (*Handler, error) {
	handler := &Handler{
		logger:    logger,
		config:    config,
		usersRepo: usersRepo,
	}

	if _, err := handler.getSecretKey(); err != nil {
		return nil, fmt.Errorf("bad secret key specified: %w", err)
	}

	return handler, nil
}
