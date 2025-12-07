package catalog

import (
	"electrotech/internal/models"
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("product not found")
)

func (r *Repo) GetProducts(p Page, filters ...FilterFunc) ([]models.Product, error) {
	products, err := r.parser.GetProducts()
	if err != nil {
		return nil, fmt.Errorf("failed get products: %s",err)
	}

	var filtered []models.Product
	for _, p := range products {
		ok := true
		for _, f := range filters {
			if !f(p) {
				ok = false
				break
			}
		}
		if ok {
			filtered = append(filtered, p)
		}
	}

	return filtered, nil
}

type FilterFunc func(p models.Product) bool

func (r *Repo) GetProduct(id string) (models.Product, error) {
	products, err := r.parser.GetProducts()
	if err != nil {
		return models.Product{}, err
	}
	for _, p := range products {
		if p.Id == id {
			return p, nil
		}
	}

	return models.Product{}, ErrNotFound
}

func (r *Repo) GetProductPrice(id string) (float32, error) {
	product, err := r.GetProduct(id)
	if err != nil {
		return 0, err
	}
	return product.Price, nil
}

func (r *Repo) GetProductName(id string) (string, error) {
	product, err := r.GetProduct(id)
	if err != nil {
		return "", err
	}
	return product.Name, nil
}
