package order

import (
	"context"
	"electrotech/internal/models"
	"fmt"
	"sync"

	"charm.land/log/v2"
)

type EmailService interface {
	SendInfo(content []byte, subject string) error
}
type OrderRepository interface {
	New(ctx context.Context, user *models.User, products []models.OrderProduct) (*models.Order, error)
	GetOrders(ctx context.Context, userID int64) ([]*models.Order, error)
}
type CatalogRepository interface {
	GetProduct(id string) (models.Product, error)
}
type UserRepository interface {
	ByID(ctx context.Context, id int64) (*models.User, error)
}

type Service struct {
	logger       *log.Logger
	repo         OrderRepository
	catalog      CatalogRepository
	userRepo     UserRepository
	emailService EmailService
}

func New(
	logger *log.Logger,
	repo OrderRepository,
	catalog CatalogRepository,
	userRepo UserRepository,
	emailService EmailService,
) *Service {
	return &Service{
		logger:       logger,
		repo:         repo,
		catalog:      catalog,
		userRepo:     userRepo,
		emailService: emailService,
	}
}

func (s *Service) GetOrders(ctx context.Context, userID int64) ([]*models.Order, error) {
	orders, err := s.repo.GetOrders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get orders from repo: %w", err)
	}

	var wg sync.WaitGroup

	for _, o := range orders {
		for _, p := range o.OrderProducts {
			wg.Go(func() {
				img, err := s.GetImagePath(ctx, p.ProductID)
				if err != nil {
					s.logger.Warn("Failed getting product, skipping", "error", err, "productID", p.ProductID)

					return
				}

				p.ImagePath = img
			})
		}
	}

	return orders, nil
}

func (s *Service) GetImagePath(ctx context.Context, productID string) (string, error) {
	product, err := s.catalog.GetProduct(productID)
	if err != nil {
		return "", fmt.Errorf("get product from repo: %w", err)
	}

	return product.ImagePath, nil
}

func (s *Service) CreateOrder(ctx context.Context, userID int64, items []models.OrderItem) (int64, error) {
	user, err := s.userRepo.ByID(ctx, userID)
	if err != nil {
		return -1, fmt.Errorf("get user by id: %w", err)
	}

	products, err := s.ItemsToProducts(ctx, items)
	if err != nil {
		return -1, fmt.Errorf("get products: %w", err)
	}

	order, err := s.repo.New(ctx, user, products)
	if err != nil {
		return -1, fmt.Errorf("insert order: %w", err)
	}

	go s.sendEmail(*order)

	return order.ID, nil
}

func (s *Service) ItemsToProducts(ctx context.Context, items []models.OrderItem) ([]models.OrderProduct, error) {
	mapped := make([]models.OrderProduct, 0, len(items))

	for _, p := range items {
		product, err := s.catalog.GetProduct(p.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %q: %w", p.ProductID, err)
		}

		mapped = append(mapped,
			models.OrderProduct{ //nolint:exhaustruct_v5
				ProductName:  product.Name,
				Quantity:     p.Quantity,
				ProductPrice: float64(product.Price),
				ProductID:    p.ProductID,
			},
		)
	}

	return mapped, nil
}
