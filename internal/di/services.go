package di

import (
	"electrotech/internal/config"
	"electrotech/internal/email"
	"electrotech/internal/repository/catalog"
	"electrotech/internal/repository/orders"
	"electrotech/internal/repository/users"
	"electrotech/internal/services/feedback"
	"electrotech/internal/services/health"
	"electrotech/internal/services/order"
	"electrotech/internal/services/user"
	"fmt"

	"charm.land/log/v2"
	"gorm.io/gorm"
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

func newHealthService(
	logger *log.Logger,
	db *gorm.DB,
) *health.Service {
	return health.New(
		health.WithLogger(logger),
		health.WithDB(db),
	)
}

func newFeedbackService(
	logger *log.Logger,
	emailService *email.Service,
	db *gorm.DB,
) *feedback.Service {
	return feedback.New(logger, emailService, db)
}
