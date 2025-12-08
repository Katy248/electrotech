package parser

import (
	"electrotech/internal/models"
	"fmt"

	"github.com/charmbracelet/log"
)

func mapProducts(offers *offersModel, imports *importsModel) ([]models.Product, error) {
	products := []models.Product{}
	for _, xmlProduct := range imports.Catalog.Products {
		p := models.Product{
			Id:            xmlProduct.Id,
			Name:          xmlProduct.Name,
			ArticleNumber: xmlProduct.ArticleNumber,
			Description:   xmlProduct.Description,
			ImagePath:     xmlProduct.Image,
		}
		manufacturer, err := getManufacturer(xmlProduct, imports)
		if err != nil {
			return nil, fmt.Errorf("failed get manufacturer for product: %s", err)
		}
		p.Manufacturer = manufacturer

		price, currency, currencySym, err := getPrice(xmlProduct, offers)
		if err != nil {
			return products, fmt.Errorf("failed get price for product: %s", err)
		}
		p.Price = price
		p.Currency = currency
		p.CurrencySym = currencySym

		count, err := getCount(xmlProduct, offers)
		if err != nil {
			return products, fmt.Errorf("failed get price for product: %s", err)
		}
		p.Count = count

		products = append(products, p)
	}
	return products, nil
}

func getPrice(p product, offers *offersModel) (price float32, currency string, currencySymbol string, err error) {
	o, err := getOffer(p, offers)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed get offer for product: %s", err)
	}

	if len(o.Prices) == 0 {
		return 0, "", "", fmt.Errorf("there is no prices specified for offer")
	}

	currency, currencySym := getCurrency(o.Prices[0])
	return o.Prices[0].Value, currency, currencySym, nil
}

func getCurrency(p price) (currency string, symbol string) {
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
		log.Warn("Unknown currency, fallback to default shekel symbol", "currency", p.Currency, "fallback to", models.CurrencySymbolILS)
		return models.CurrencyILS, models.CurrencySymbolILS
	}
}

func getOffer(p product, off *offersModel) (*offer, error) {
	for _, o := range off.Package.Offers {
		if o.Id == p.Id {
			return &o, nil
		}

	}
	return nil, fmt.Errorf("there is no offer for product id = '%s'", p.Id)
}

func getCount(p product, offers *offersModel) (float32, error) {
	o, err := getOffer(p, offers)
	if err != nil {
		return 0, fmt.Errorf("failed get offer for product: %s", err)
	}
	return o.Count, nil
}

func getManufacturer(p product, imports *importsModel) (string, error) {
	if len(p.GroupIds) == 0 {
		return "", fmt.Errorf("there is no group specified for product")
	}
	id := p.GroupIds[0]
	group, err := imports.getGroup(id)
	if err != nil {
		return "", fmt.Errorf("failed get group: %s", err)
	}
	return group.Name, nil
}

func (i *importsModel) getGroup(id string) (*group, error) {
	for _, g := range i.Classifier.Groups {
		if g.Id == id {
			return &g, nil
		}
	}

	return nil, fmt.Errorf("there is no group with id = '%s'", id)
}
