package parser

import (
	"electrotech/internal/models"
	"fmt"
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

		price, err := getPrice(xmlProduct, offers)
		if err != nil {
			return products, fmt.Errorf("failed get price for product: %s", err)
		}
		p.Price = price

		count, err := getCount(xmlProduct, offers)
		if err != nil {
			return products, fmt.Errorf("failed get price for product: %s", err)
		}
		p.Count = count

		products = append(products, p)
	}
	return products, nil
}

func getPrice(p product, offers *offersModel) (float32, error) {
	o, err := getOffer(p, offers)
	if err != nil {
		return 0, fmt.Errorf("failed get offer for product: %s", err)
	}

	if len(o.Prices) == 0 {
		return 0, fmt.Errorf("there is no prices specified for offer")
	}

	return o.Prices[0].Value, nil
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
