package catalog

import (
	"context"
	"electrotech/internal/models"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound = errors.New("product not found")
)

const PageSize = 20

type Products struct {
	Products []models.Product
	Page     int
	Pages    int
	Total    int
}

func (r *Repo) GetProductsNew(ctx context.Context, p Page, filters ...FilterFunc) (*Products, error) {
	products, err := r.parser.GetProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("get products from parser: %w", err)
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
		return nil, ErrNotFound
	}

	total := len(filtered)

	pages := GetPages(len(filtered), PageSize)

	if int(p*PageSize) > len(filtered)-1 {
		return nil, ErrNotFound
	} else {
		filtered = filtered[int(p*PageSize):]
	}

	return &Products{
		Products: takeFirst(filtered, PageSize),
		Pages:    pages,
		Page:     int(p),
		Total:    total,
	}, nil
}

func GetPages(productsCount int, pageSize int) int {
	pages := productsCount / pageSize
	if productsCount%pageSize != 0 {
		pages++
	}

	return pages
}

func takeFirst(products []models.Product, nFirst int) []models.Product {
	filtered := []models.Product{}

	for i, p := range products {
		if i == nFirst {
			break
		}

		filtered = append(filtered, p)
	}

	return filtered
}

type FilterFunc func(p models.Product) bool

func QueryFilter(query string) FilterFunc {
	return func(p models.Product) bool {
		return strings.Contains(p.Name, query) ||
			strings.Contains(p.Description, query) ||
			strings.Contains(p.Id, query) ||
			strings.Contains(p.ArticleNumber, query) ||
			strings.Contains(p.ImagePath, query) ||
			strings.Contains(p.Manufacturer, query)
	}
}
func OnlyAvailableFilter() FilterFunc {
	return func(p models.Product) bool {
		return p.Count > 0
	}
}

type result struct {
	product models.Product
	err     error
}

func (r *Repo) GetProduct(ctx context.Context, id string) (models.Product, error) {
	resChan := make(chan result, 1)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func(ctx context.Context, resChan chan result) {
		products, err := r.parser.GetProducts(ctx)
		if err != nil {
			resChan <- result{err: fmt.Errorf("parse products: %w", err)} //nolint:exhaustruct_v5

			return
		}

		for _, p := range products {
			if p.Id == id {
				resChan <- result{product: p} //nolint:exhaustruct_v5

				return
			}
		}

		resChan <- result{err: ErrNotFound} //nolint:exhaustruct_v5
	}(ctx, resChan)

	select {
	case res := <-resChan:
		if res.err != nil {
			return models.Product{}, res.err
		}

		return res.product, nil
	case <-ctx.Done():
		return models.Product{}, fmt.Errorf("context canceled: %w", ctx.Err())
	}
}

func (r *Repo) GetProductPrice(ctx context.Context, id string) (float32, error) {
	product, err := r.GetProduct(ctx, id)
	if err != nil {
		return 0, err
	}

	return product.Price, nil
}

func (r *Repo) GetProductName(ctx context.Context, id string) (string, error) {
	product, err := r.GetProduct(ctx, id)
	if err != nil {
		return "", err
	}

	return product.Name, nil
}
