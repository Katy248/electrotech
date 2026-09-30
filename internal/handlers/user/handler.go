package user

import (
	"electrotech/internal/repository/users"

	"charm.land/log/v2"
)

type Handler struct {
	logger    *log.Logger
	usersRepo *users.Repo
}

func NewHandler(logger *log.Logger, usersRepo *users.Repo) *Handler {
	return &Handler{
		logger:    logger,
		usersRepo: usersRepo,
	}
}
