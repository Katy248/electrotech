package parser

import (
	"electrotech/internal/models"
	"errors"
	"fmt"

	"charm.land/log/v2"
)

func mapProducts(offers *offersModel, imports *importsModel) ([]models.Product, error) {
	products := []models.Product{}

	for _, xmlProduct := range imports.Catalog.Products {
		product := models.Product{
			Id:            xmlProduct.Id,
			Name:          xmlProduct.Name,
			ArticleNumber: xmlProduct.ArticleNumber,
			Description:   xmlProduct.Description,
			ImagePath:     xmlProduct.Image,
			Count:         0,  // default
			Currency:      "", // default
			CurrencySym:   "", // default
			Price:         0,  // default
			Manufacturer:  "",
			Category:      models.NilCategory,
		}

		manufacturer, err := getManufacturer(xmlProduct, imports)
		if err != nil {
			return nil, fmt.Errorf("failed get manufacturer for product: %w", err)
		}

		product.Manufacturer = manufacturer

		category, err := getCategory(xmlProduct, imports)
		if err != nil {
			return nil, fmt.Errorf("failed get category for product: %w", err)
		}

		product.Category = category

		price, currency, currencySym, err := getPrice(xmlProduct, offers)
		if err != nil {
			return products, fmt.Errorf("failed get price for product: %w", err)
		}

		product.Price = price
		product.Currency = currency
		product.CurrencySym = currencySym

		count, err := getCount(xmlProduct, offers)
		if err != nil {
			return products, fmt.Errorf("failed get price for product: %w", err)
		}

		product.Count = count

		products = append(products, product)
	}

	return products, nil
}

var ErrEmptyCategoryID = errors.New("category ID is empty")
var ErrCategoryNotFound = errors.New("category not found")

func getCategory(p product, imports *importsModel) (models.Category, error) {
	if p.CategoryId == "" {
		return models.Category{}, ErrEmptyCategoryID
	}

	for _, i := range imports.Classifier.Categories {
		if i.ID == p.CategoryId {
			return models.Category{
				Name: i.Name,
				Id:   i.ID,
			}, nil
		}
	}

	return models.Category{}, ErrCategoryNotFound
}

var ErrNoPrices = errors.New("there is no prices specified for offer")

func getPrice(p product, offers *offersModel) (float32, string, string, error) {
	o, err := getOffer(p, offers)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed get offer for product: %w", err)
	}

	if len(o.Prices) == 0 {
		return 0, "", "", ErrNoPrices
	}

	currency, currencySym := getCurrency(o.Prices[0])

	return o.Prices[0].Value, currency, currencySym, nil
}

func getCurrency(p price) (string, string) {
	switch p.Currency {
	case "руб":
		return models.CurrencyRUB, models.CurrencySymbolRUB
	case "EUR":
		return models.CurrencyEUR, models.CurrencySymbolEUR
	case "USD":
		return models.CurrencyUSD, models.CurrencySymbolUSD
	case "ILS":
		return models.CurrencyILS, models.CurrencySymbolILS
	default:
		log.Warn(
			"Unknown currency, fallback to default shekel symbol",
			"currency", p.Currency,
			"fallback to", models.CurrencySymbolILS,
		)

		return models.CurrencyILS, models.CurrencySymbolILS
	}
}

var ErrNoOffer = errors.New("there is no offer for product")

func getOffer(p product, off *offersModel) (*offer, error) {
	for _, o := range off.Package.Offers {
		if o.Id == p.Id {
			return &o, nil
		}
	}

	return nil, ErrNoOffer
}

func getCount(p product, offers *offersModel) (float32, error) {
	o, err := getOffer(p, offers)
	if err != nil {
		return 0, fmt.Errorf("failed get offer for product: %w", err)
	}

	return o.Count, nil
}

var ErrNoGroup = errors.New("there is no group specified for product")

func getManufacturer(p product, imports *importsModel) (string, error) {
	if len(p.GroupIds) == 0 {
		return "", ErrNoGroup
	}

	id := p.GroupIds[0]

	group, err := imports.getGroup(id)
	if err != nil {
		return "", fmt.Errorf("failed get group: %w", err)
	}

	return group.Name, nil
}

func (i *importsModel) getGroup(id string) (*group, error) {
	for _, g := range i.Classifier.Groups {
		if g.Id == id {
			return &g, nil
		}
	}

	return nil, ErrNoGroup
}
