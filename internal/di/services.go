package di

import (
	"electrotech/internal/config"
	"electrotech/internal/repository/users"
	"electrotech/internal/services/user"
	"fmt"

	"charm.land/log/v2"
)

func newUserService(
	logger *log.Logger,
	repo *users.Repo,
	config *config.Config,
) (*user.Service, error) {
	s, err := user.New(logger, repo, &config.Auth)
	if err != nil {
		return nil, fmt.Errorf("new user service: %w", err)
	}

	return s, nil
}
