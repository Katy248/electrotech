package handlers

import (
	"context"
	"electrotech/internal/gen/handler_gen"
	"electrotech/internal/models"
	"electrotech/internal/services/health"

	"charm.land/log/v2"
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
	logger        *log.Logger
	orderService  OrderService
	userService   UserService
	healthService *health.Service
}

var _ handler_gen.Handler = (*Handler)(nil)

func New(
	logger *log.Logger,
	orderService OrderService,
	userService UserService,
	healthService *health.Service,
) *Handler {
	return &Handler{
		logger:        logger,
		orderService:  orderService,
		userService:   userService,
		healthService: healthService,
	}
}

// GetHealth implements [handler.Handler].
//
//nolint:ireturn
func (h *Handler) GetHealth(
	ctx context.Context,
) (handler_gen.GetHealthRes, error) {
	status := h.healthService.GetHealth(ctx)

	deps := handler_gen.HealthDependencies{}
	for name, status := range status.Dependencies {
		deps[name] = mapHealthStatus(status)
	}

	if status.Status() == health.StatusOK {
		return &handler_gen.GetHealthOK{
			Status:       mapHealthStatus(status.Status()),
			Dependencies: deps,
		}, nil
	}

	return &handler_gen.GetHealthBadGateway{
		Status:       mapHealthStatus(status.Status()),
		Dependencies: deps,
	}, nil
}

func mapHealthStatus(s health.Status) handler_gen.HealthStatus {
	if s == health.StatusOK {
		return handler_gen.HealthStatusOk
	}

	return handler_gen.HealthStatusUnavailable
}
