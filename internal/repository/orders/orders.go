package orders

import (
	"electrotech/internal/models"
	"electrotech/storage"
	"errors"
	"fmt"
	"time"

	"charm.land/log/v2"
)

var ErrUserIsNil = errors.New("user is nil")

func InsertNew(o *models.Order) error {
	err := storage.DB.Create(o).Error

	return err
}

func New(user *models.User, products []models.OrderProduct) (*models.Order, error) {
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

	err = storage.DB.Create(o).Error
	if err != nil {
		return nil, fmt.Errorf("failed insert order: %w", err)
	}

	for _, p := range products {
		o.AddProduct(p)

		err := storage.DB.Save(&p).Error
		if err != nil {
			return nil, fmt.Errorf("failed save product: %w", err)
		}
	}

	err = storage.DB.Save(o).Error
	if err != nil {
		return nil, fmt.Errorf("failed save order: %w", err)
	}

	return o, nil
}

var getOrdersQuery = `
SELECT id, user_id, creation_date
FROM orders
WHERE user_id = ?`

func GetOrders(userID int64) ([]*models.Order, error) {
	var orders []*models.Order

	err := storage.DB.Raw(getOrdersQuery, userID).Find(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("failed get orders: %w", err)
	}

	for _, o := range orders {
		err = storage.DB.Where("order_id = ?", o.ID).Find(&o.OrderProducts).Error
		log.Info("Order", "orderID", o.ID, "products", o.OrderProducts)

		if err != nil {
			return nil, fmt.Errorf("failed get products: %w", err)
		}
	}

	return orders, err
}
