package models

type Product struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	ImagePath     string  `json:"imagePath"`
	Price         float32 `json:"price"`
	ArticleNumber string  `json:"articleNumber"`
	Count         float32 `json:"count"`
	Manufacturer  string  `json:"manufacturer"`
	CurrencySym   string  `json:"currencySym"`
}

const (
	CurrencyRUB = "₽"
	CurrencyUSD = "$"
	CurrencyEUR = "€"
	CurrencyILS = "₪" // Israeli New Sheqel
)
