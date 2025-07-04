package models

type Product struct {
	Id            string
	Name          string
	Description   string
	ImagePath     string
	Price         float32
	ArticleNumber string
	Count         int
	Manufacturer  string
	Parameters    ParametersList
}

type ParametersList = map[string]string
