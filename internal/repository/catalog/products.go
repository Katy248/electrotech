package catalog

import (
	"electrotech/internal/models"
	"errors"
	"log"
	"slices"
)

var (
	ErrNotFound = errors.New("product not found")
)

func ListFilter(parameterName string, values []string) FilterFunc {
	return func(p models.Product) bool {

		parameter, err := p.GetParameter(parameterName)
		if err != nil {
			log.Printf("Error getting parameter: %v", err)
			return false
		}

		if parameter.Type != models.ParameterTypeList {
			log.Printf("Wrong filter with type %q for parameter with type %q", models.ParameterTypeList, parameter.Type)
			return false
		}

		if parameter.StringValue != "" {
			return slices.Contains(values, parameter.StringValue)
		}

		for _, v := range values {
			if !slices.Contains(parameter.SliceValue, v) {
				return false
			}
		}

		return true
	}
}
func RangeFilter(parameterName string, min, max float64) FilterFunc {
	return func(p models.Product) bool {

		parameter, err := p.GetParameter(parameterName)
		if err != nil {
			log.Printf("Error getting parameter: %v", err)
			return false
		}

		if parameter.Type != models.ParameterTypeNumber {
			log.Printf("Wrong filter with type %q for parameter with type %q", models.ParameterTypeNumber, parameter.Type)
			return false
		}

		var inMinRange, inMaxRange bool

		inMinRange = parameter.NumberValue >= min
		if max > 0 {
			inMaxRange = parameter.NumberValue <= max
		} else {
			inMaxRange = true
		}
		log.Printf("Parameter %q has value %f, min - %f, max- %f, result - %v", parameterName, parameter.NumberValue, min, max, inMinRange && inMaxRange)
		return inMinRange && inMaxRange
	}
}

func (r *Repo) GetProducts(p Page, filters ...FilterFunc) ([]models.Product, error) {
	products, err := r.parser.GetProducts()
	if err != nil {
		return nil, err
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
