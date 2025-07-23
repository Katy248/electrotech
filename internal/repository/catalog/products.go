package catalog

import (
	"electrotech/internal/models"
	"errors"
)

var (
	ErrNotFound = errors.New("product not found")
)

func (r *Repo) GetProducts(p Page) ([]models.Product, error) {
	return r.parser.GetProducts()
}

type Filter struct {
}

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

func (r *Repo) GetProductsFilter(p Page) ([]models.Product, error) {
	return nil, ErrNotImplemented
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
