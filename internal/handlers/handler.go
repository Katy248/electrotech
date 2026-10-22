package handlers

import (
	"context"
	"electrotech/internal/gen/handler_gen"
	"electrotech/internal/models"

	"charm.land/log/v2"
	"gorm.io/gorm"
)

type OrderService interface {
	GetOrders(ctx context.Context, userID int64) ([]*models.Order, error)
	CreateOrder(ctx context.Context, userID int64, items []models.OrderItem) (int64, error)
}
type UserService interface {
	UpdateEmail(ctx context.Context, id int64, email string) error
	UpdatePassword(ctx context.Context, id int64, old, new string) error
	UpdatePhone(ctx context.Context, id int64, phone string) error
	GetUser(ctx context.Context, id int64) (*models.User, error)
	Update(ctx context.Context, id int64, update *models.UserUpdate) error
}

type Handler struct {
	db           *gorm.DB
	logger       *log.Logger
	orderService OrderService
	userService  UserService
}

var _ handler_gen.Handler = (*Handler)(nil)

func New(
	db *gorm.DB,
	logger *log.Logger,
	orderService OrderService,
	userService UserService,
) *Handler {
	return &Handler{
		db:           db,
		logger:       logger,
		orderService: orderService,
		userService:  userService,
	}
}

// GetHealth implements [handler.Handler].
//
//nolint:ireturn
func (h *Handler) GetHealth(
	ctx context.Context,
) (handler_gen.GetHealthRes, error) {
	sqlDB, err := h.db.DB()
	if err != nil {
		h.logger.Error("Get SQL database from ORM", "error", err)

		return &handler_gen.GetHealthBadGateway{
			Status: handler_gen.HealthStatusUnavailable,
			Dependencies: handler_gen.HealthDependencies{
				Db: handler_gen.HealthStatusUnavailable,
			},
		}, nil
	}

	if err := sqlDB.Ping(); err != nil {
		h.logger.Error("Ping database", "error", err)

		return &handler_gen.GetHealthBadGateway{
			Status: handler_gen.HealthStatusUnavailable,
			Dependencies: handler_gen.HealthDependencies{
				Db: handler_gen.HealthStatusUnavailable,
			},
		}, nil
	}

	return &handler_gen.GetHealthBadGateway{
		Status: handler_gen.HealthStatusOk,
		Dependencies: handler_gen.HealthDependencies{
			Db: handler_gen.HealthStatusOk,
		},
	}, nil
}
