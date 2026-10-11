package di

import (
	"electrotech/internal/config"
	"electrotech/internal/repository/users"
	"electrotech/internal/services/user"

	"charm.land/log/v2"
)

func newUserService(
	logger *log.Logger,
	repo *users.Repo,
	config *config.Config,
) *user.Service {
	return user.New(logger, repo, &config.Auth)
}
