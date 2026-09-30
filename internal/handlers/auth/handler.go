package auth

import (
	"electrotech/internal/config"
	"electrotech/internal/repository/users"
	"fmt"

	"charm.land/log/v2"
)

type Handler struct {
	logger    *log.Logger
	config    *config.AuthConfig
	usersRepo *users.Repo
}

func NewHandler(logger *log.Logger, config *config.AuthConfig, usersRepo *users.Repo) (*Handler, error) {
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
