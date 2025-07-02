package parser

import "encoding/xml"

// imports

type baseUnit struct {
	XMLName                   xml.Name `xml:"БазоваяЕдиница"`
	Code                      int      `xml:"Код,attr"`
	FullName                  string   `xml:"ПолноеНаименование,attr"`
	InternationalAbbreviation string   `xml:"МеждународноеСокращение,attr"`
}

type product struct {
	Id             string          `xml:"Ид"`
	ArticleNumber  string          `xml:"Артикул"`
	Code           string          `xml:"Код"`
	Name           string          `xml:"Наименование"`
	GroupIds       []string        `xml:"Группы>Ид"`
	CategoryId     string          `xml:"Категория"`
	Description    string          `xml:"Описание"`
	Country        string          `xml:"Страна"`
	Image          string          `xml:"Картинка"`
	PropertyValues []propertyValue `xml:"ЗначенияСвойств>ЗначенияСвойства"`
	ProductUnit    baseUnit
}
type group struct {
	Id   string `xml:"Ид"`
	Name string `xml:"Наименование"`
}
type classifier struct {
	Groups []group `xml:"Группы>Группа"`
}

type propertyValue struct {
	Id    string `xml:"Ид"`
	Value string `xml:"Значение"`
}
type catalog struct {
	ContainsOnlyChanges bool      `xml:"СодержитТолькоИзменения,attr"`
	Products            []product `xml:"Товары>Товар"`
}
type importsModel struct {
	Catalog    catalog    `xml:"Каталог"`
	Classifier classifier `xml:"Классификатор"`
}

// offers

type price struct {
	Presentation string  `xml:"Представление"`
	PriceTypeId  string  `xml:"ИдТипаЦены"`
	Value        float32 `xml:"ЦенаЗаЕдиницу"`
	Currency     string  `xml:"Валюта"`
	Unit         string  `xml:"Единица"`
	Ratio        int     `xml:"Коэффициент"`
}

type offer struct {
	Id            string `xml:"Ид"`
	Name          string `xml:"Наименование"`
	Unit          baseUnit
	ArticleNumber string  `xml:"Артикул"`
	Count         int     `xml:"Количество"`
	Prices        []price `xml:"Цены>Цена"`
}
type priceType struct {
	Id       string `xml:"Ид"`
	Name     string `xml:"Наименование"`
	Currency string `xml:"Валюта"`
	Tax      struct {
		Name     string `xml:"Наименование"`
		Included bool   `xml:"УчтеноВСумме"`
	} `xml:"Налог"`
}

type offersPackage struct {
	XMLName             xml.Name    `xml:"ПакетПредложений"`
	Offers              []offer     `xml:"Предложения>Предложение"`
	ContainsOnlyChanges bool        `xml:"СодержитТолькоИзменения,attr"`
	PriceTypes          []priceType `xml:"ТипыЦен>ТипЦены"`
}

type offersModel struct {
	Package offersPackage `xml:"ПакетПредложений"`
}
