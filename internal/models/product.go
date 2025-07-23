package models

type Product struct {
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	ImagePath     string         `json:"imagePath"`
	Price         float32        `json:"price"`
	ArticleNumber string         `json:"articleNumber"`
	Count         int            `json:"count"`
	Manufacturer  string         `json:"manufacturer"`
	Parameters    ParametersList `json:"parameters"`
}

type ParametersList = map[string]string
