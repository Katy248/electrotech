package models

import "fmt"

type Product struct {
	Id            string             `json:"id"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	ImagePath     string             `json:"imagePath"`
	Price         float32            `json:"price"`
	ArticleNumber string             `json:"articleNumber"`
	Count         int                `json:"count"`
	Manufacturer  string             `json:"manufacturer"`
	Parameters    []ProductParameter `json:"parameters"`
}

type ProductParameter struct {
	Name        string        `json:"name"`
	Type        ParameterType `json:"type"`
	StringValue string        `json:"stringValue"`
	NumberValue float64       `json:"numberValue"`
	SliceValue  []string      `json:"sliceValue"`
}

func (p *Product) GetParameter(name string) (*ProductParameter, error) {
	for _, param := range p.Parameters {
		if param.Name == name {
			return &param, nil
		}
	}
	return nil, fmt.Errorf("parameter %q not found", name)
}
