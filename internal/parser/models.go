package parser

import "encoding/xml"

// imports

type BaseUnit struct {
	XMLName                   xml.Name `xml:"БазоваяЕдиница"`
	Code                      int      `xml:"Код,attr"`
	FullName                  string   `xml:"ПолноеНаименование,attr"`
	InternationalAbbreviation string   `xml:"МеждународноеСокращение,attr"`
}

type Product struct {
	Id             string `xml:"Ид"`
	ProductUnit    BaseUnit
	ArticleNumber  string          `xml:"Артикул"`
	Code           string          `xml:"Код"`
	Name           string          `xml:"Наименование"`
	GroupIds       []string        `xml:"Группы>Ид"`
	CategoryId     string          `xml:"Категория"`
	Description    string          `xml:"Описание"`
	Country        string          `xml:"Страна"`
	Image          string          `xml:"Картинка"`
	PropertyValues []PropertyValue `xml:"ЗначенияСвойств>ЗначенияСвойства"`
}
type PropertyValue struct {
	Id    string `xml:"Ид"`
	Value string `xml:"Значение"`
}
type Catalog struct {
	ContainsOnlyChanges bool      `xml:"СодержитТолькоИзменения,attr"`
	Products            []Product `xml:"Товары>Товар"`
}
type ImportsModel struct {
	Catalog Catalog `xml:"Каталог"`
}

// offers

type Price struct {
	Presentation string  `xml:"Представление"`
	PriceTypeId  string  `xml:"ИдТипаЦены"`
	Value        float64 `xml:"ЦенаЗаЕдиницу"`
	Currency     string  `xml:"Валюта"`
	Unit         string  `xml:"Единица"`
	Ratio        int     `xml:"Коэффициент"`
}

type Offer struct {
	Id            string `xml:"Ид"`
	Name          string `xml:"Наименование"`
	Unit          BaseUnit
	ArticleNumber string  `xml:"Артикул"`
	Count         int     `xml:"Количество"`
	Prices        []Price `xml:"Цены>Цена"`
}
type PriceType struct {
	Id       string `xml:"Ид"`
	Name     string `xml:"Наименование"`
	Currency string `xml:"Валюта"`
	Tax      struct {
		Name     string `xml:"Наименование"`
		Included bool   `xml:"УчтеноВСумме"`
	} `xml:"Налог"`
}

type OffersPackage struct {
	XMLName             xml.Name    `xml:"ПакетПредложений"`
	Offers              []Offer     `xml:"Предложения>Предложение"`
	ContainsOnlyChanges bool        `xml:"СодержитТолькоИзменения,attr"`
	PriceTypes          []PriceType `xml:"ТипыЦен>ТипЦены"`
}

type OffersModel struct {
	Package OffersPackage `xml:"ПакетПредложений"`
}
