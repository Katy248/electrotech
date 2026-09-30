package orders

import (
	"electrotech/internal/models"
	"errors"
	"fmt"
	"time"

	"charm.land/log/v2"
	"gorm.io/gorm"
)

var ErrUserIsNil = errors.New("user is nil")

type Repo struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewRepo(db *gorm.DB, logger *log.Logger) *Repo {
	return &Repo{db: db, logger: logger}
}

func (r *Repo) InsertNew(o *models.Order) error {
	err := r.db.Create(o).Error

	return err
}

func (r *Repo) New(user *models.User, products []models.OrderProduct) (*models.Order, error) {
	if user == nil {
		return nil, ErrUserIsNil
	}

	o := &models.Order{ //nolint:exhaustruct_v5
		CreationDate: time.Now(),
	}

	err := o.SetUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed set user: %w", err)
	}

	err = r.db.Create(o).Error
	if err != nil {
		return nil, fmt.Errorf("failed insert order: %w", err)
	}

	for _, p := range products {
		o.AddProduct(p)

		err := r.db.Save(&p).Error
		if err != nil {
			return nil, fmt.Errorf("failed save product: %w", err)
		}
	}

	err = r.db.Save(o).Error
	if err != nil {
		return nil, fmt.Errorf("failed save order: %w", err)
	}

	return o, nil
}

const getOrdersQuery = `
SELECT id, user_id, creation_date
FROM orders
WHERE user_id = ?`

func (r *Repo) GetOrders(userID int64) ([]*models.Order, error) {
	var orders []*models.Order

	err := r.db.Raw(getOrdersQuery, userID).Find(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("failed get orders: %w", err)
	}

	for _, o := range orders {
		err = r.db.Where("order_id = ?", o.ID).Find(&o.OrderProducts).Error
		r.logger.Info("Order", "orderID", o.ID, "products", o.OrderProducts)

		if err != nil {
			return nil, fmt.Errorf("failed get products: %w", err)
		}
	}

	return orders, err
}
