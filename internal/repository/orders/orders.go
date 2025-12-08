package orders

import (
	"electrotech/internal/models"
	"electrotech/storage"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
)

func InsertNew(o *models.Order) error {
	err := storage.DB.Create(o).Error
	return err
}

func New(user *models.User, products []models.OrderProduct) (*models.Order, error) {
	o := &models.Order{
		CreationDate: time.Now(),
	}
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	if err := o.SetUser(user); err != nil {
		return nil, fmt.Errorf("failed set user: %s", err)
	}
	err := storage.DB.Create(o).Error
	if err != nil {
		return nil, fmt.Errorf("failed insert order: %s", err)
	}

	for _, p := range products {
		o.AddProduct(p)
		err = storage.DB.Save(&p).Error
		if err != nil {
			return nil, fmt.Errorf("failed save product: %s", err)
		}
	}
	err = storage.DB.Save(o).Error
	if err != nil {
		return nil, fmt.Errorf("failed save order: %s", err)
	}

	return o, nil
}

var getOrdersQuery = `
SELECT *
FROM orders
WHERE user_id = ?`

func GetOrders(userID int64) ([]*models.Order, error) {
	var orders []*models.Order
	err := storage.DB.Raw(getOrdersQuery, userID).Find(&orders).Error
	if err != nil {
		return nil, fmt.Errorf("failed get orders: %s", err)
	}

	for _, o := range orders {
		err = storage.DB.Where("order_id = ?", o.ID).Find(&o.OrderProducts).Error
		log.Info("Order", "orderID", o.ID, "products", o.OrderProducts)
		if err != nil {
			return nil, fmt.Errorf("failed get products: %s", err)
		}
	}

	return orders, err
}
