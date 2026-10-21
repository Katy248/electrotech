package di

import (
	"electrotech/internal/config"
	"electrotech/internal/email"
	"electrotech/internal/repository/catalog"
	"electrotech/internal/repository/orders"
	"electrotech/internal/repository/users"
	"electrotech/internal/services/order"
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

func newOrderService(
	logger *log.Logger,
	orderRepo *orders.Repo,
	catalogRepo *catalog.Repo,
	userRepo *users.Repo,
	emailService *email.Service,
) *order.Service {
	return order.New(logger, orderRepo, catalogRepo, userRepo, emailService)
}
