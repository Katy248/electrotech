package catalog

import "electrotech/internal/models"

func (r *CatalogRepo) GetProducts(p Page) ([]models.Product, error) {
	return r.parser.GetProducts()
}

type Filter struct {
}

func (r *CatalogRepo) GetProductsFilter(p Page) ([]models.Product, error) {
	return nil, ErrNotImplemented
}
