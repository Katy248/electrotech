package user

import (
	"context"
	"electrotech/internal/models"

	"charm.land/log/v2"
)

type UserRepository interface {
	ByID(ctx context.Context, id int64) (*models.User, error)
	ByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
}

type Handler struct {
	logger    *log.Logger
	usersRepo UserRepository
}

func NewHandler(logger *log.Logger, usersRepo UserRepository) *Handler {
	return &Handler{
		logger:    logger,
		usersRepo: usersRepo,
	}
}
