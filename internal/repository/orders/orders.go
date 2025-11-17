package orders

import (
	"electrotech/internal/models"
	"electrotech/storage"
)

func InsertNew(o *models.Order) error {
	err := storage.DB.Create(&o).Error
	return err
}

func GetOrders(userID int64) ([]models.Order, error) {
	var orders []models.Order
	err := storage.DB.Where("user_id = ?", userID).Preload("OrderProducts").Find(&orders).Error
	return orders, err
}
