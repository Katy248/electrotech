package models

type Parameter struct {
	Name  string
	Value string
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
