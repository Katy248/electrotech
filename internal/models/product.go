package models

type ParameterType = string

const (
	ParameterTypeNumber     ParameterType = "number"
	ParameterTypeString     ParameterType = "str"
	ParameterTypeDictionary ParameterType = "dict"
)

type Parameter struct {
	Name string        `json:"name"`
	Type ParameterType `json:"type"`
	// Value string
	Values []string `json:"values,omitzero"`
}

type Product struct {
	Id            string
	Name          string
	Description   string
	ImagePath     string
	Price         float32
	ArticleNumber string
	Count         int
	Manufacturer  string
}
