package models

import (
	"fmt"
	"time"
)

type Order struct {
	ID            int64           `json:"id" gorm:"primaryKey"`
	UserID        int64           `json:"userId"`
	CreationDate  time.Time       `json:"creationDate"`
	OrderProducts []*OrderProduct `json:"products"`
	User          *User           `json:"user"`
}

func (o Order) Sum() float64 {
	sum := 0.0
	for _, p := range o.OrderProducts {
		sum += p.Sum()
	}
	return sum
}

func NewOrder() *Order {
	return &Order{
		CreationDate: time.Now(),
	}
}

func (o *Order) SetUser(u *User) error {
	if u == nil {
		return fmt.Errorf("user is nil")
	}
	if !u.CompanyData().DataFilled() {
		return fmt.Errorf("user has no required company data filled")
	}
	o.UserID = u.ID
	o.User = u
	return nil
}

// TODO: Checking duplicates
func (o *Order) AddProduct(op OrderProduct) {
	op.OrderID = o.ID
	op.Order = *o
	o.OrderProducts = append(o.OrderProducts, &op)
}

type OrderProduct struct {
	ID           int64   `json:"id"`
	OrderID      int64   `json:"orderId"`
	Order        Order   `json:"-"`
	ProductName  string  `json:"productName"`
	Quantity     int64   `json:"quantity"`
	ProductPrice float64 `json:"productPrice"`
	ProductID    string  `json:"productId"`
	ImagePath    string  `json:"imagePath"`
}

func (p OrderProduct) Sum() float64 {
	return float64(p.Quantity) * p.ProductPrice
}
