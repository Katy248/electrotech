package catalog

import (
	"electrotech/internal/models"
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("product not found")
)

const PageSize = 20

// DEPRECATED: Use GetProductsNew instead
func (r *Repo) GetProducts(p Page, filters ...FilterFunc) ([]models.Product, error) {
	products, err := r.parser.GetProducts()
	if err != nil {
		return nil, fmt.Errorf("failed get products: %s", err)
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

	if len(filtered) == 0 {
		return nil, nil
	}

	if int(p*PageSize) > len(filtered) {
		return nil, nil
	} else {
		filtered = filtered[int(p*PageSize):]
	}

	return takeFirst(filtered, PageSize), nil
}
func (r *Repo) GetProductsNew(p Page, filters ...FilterFunc) ([]models.Product, int, error) {
	products, err := r.parser.GetProducts()
	if err != nil {
		return nil, 0, fmt.Errorf("failed get products: %s", err)
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

	if len(filtered) == 0 {
		return nil, 0, nil
	}

	pages := len(filtered) / PageSize
	if len(filtered)%PageSize != 0 {
		pages++
	}

	if int(p*PageSize) > len(filtered) {
		return nil, 0, nil
	} else {
		filtered = filtered[int(p*PageSize):]
	}

	return takeFirst(filtered, PageSize), pages, nil
}

func takeFirst(products []models.Product, n int) []models.Product {
	filtered := []models.Product{}
	for i := 0; i < n; i++ {
		filtered = append(filtered, products[i])
	}
	return filtered
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
