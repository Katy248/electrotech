package models

import (
	"errors"
	"time"
)

type Order struct {
	ID            int64           `gorm:"primaryKey"   json:"id"`
	UserID        int64           `json:"userId"`
	CreationDate  time.Time       `json:"creationDate"`
	OrderProducts []*OrderProduct `json:"products"`
	User          *User           `json:"user"`
}

func NewOrder() *Order {
	return &Order{
		CreationDate: time.Now(),
	}
}

func (o Order) Sum() float64 {
	sum := 0.0
	for _, p := range o.OrderProducts {
		sum += p.Sum()
	}

	return sum
}

var ErrUserIsNil = errors.New("user is nil")

var ErrCompanyDataNotFilled = errors.New("user has no required company data filled")

func (o *Order) SetUser(u *User) error {
	if u == nil {
		return ErrUserIsNil
	}

	if !u.CompanyData().DataFilled() {
		return ErrCompanyDataNotFilled
	}

	o.UserID = u.ID
	o.User = u

	return nil
}

// AddProduct adds a product to the order.
//
// TODO: Checking duplicates.
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
