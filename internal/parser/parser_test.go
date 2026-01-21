package parser

import (
	"os"
	"testing"
)

func TestNewParser(t *testing.T) {

	dir, _ := os.Getwd()
	t.Logf("Current directory: %s", dir)

	_, err := NewParser("../../example")
	if err != nil {
		t.Errorf("Failed create new parser: %s", err)
	}

	_, expectedErr := NewParser("./not_exists")
	if expectedErr == nil {
		t.Errorf("Expected error when directory not exists")
	}
}

func TestParseImportsData(t *testing.T) {
	result, err := parseImportsData(importsData)

	if err != nil {
		t.Error(err)
	}

	if result.Catalog.ContainsOnlyChanges != true {
		t.Error("ContainsOnlyChanges (xml attribute 'СодержитТолькоИзменения') failed to parse (should be true)")
	}

	if len(result.Catalog.Products) <= 0 {
		t.Fatal("Products failed to parse (should be not zero items)")
	}

	firstProduct := result.Catalog.Products[0]

	if firstProduct.Id != "RandomId" {
		t.Error("Id failed to parse (should be 'RandomId')")
	}

	if len(result.Classifier.Categories) < 1 {
		t.Error("Categories failed to parse (should be not zero items)")
	}

	firstCategory := result.Classifier.Categories[0]

	expectedCategoryId := "b0627bbe-ceb3-11eb-a53c-f14c52115a4f"
	expectedCategoryName := "<Без категории>"
	if firstCategory.ID != expectedCategoryId {
		t.Errorf("Category.Id failed to parse (should be '%s'), but was '%s'", expectedCategoryId, firstCategory.ID)
	}
	if firstCategory.Name != expectedCategoryName {
		t.Errorf("Category.Name failed to parse (should be '%s') but was '%s'", expectedCategoryName, firstCategory.Name)
	}

	if len(firstProduct.GroupIds) <= 0 {
		t.Fatal("Groups failed to parse (should be not zero items)")
	}
	if firstProduct.GroupIds[0] != "GroupId" {
		t.Error("Groups failed to parse (should be 'GroupId')")
	}

	for _, prop := range result.Classifier.Properties {
		if prop.Type == propertyTypeHandbook {
			if len(prop.Variants) <= 0 {
				t.Errorf("Property %s of type handbook but has no variants", prop.Name)
			}
		}
	}

}

func TestParseOffersData(t *testing.T) {
	data := offersData
	result, err := parseOffersData(data)
	if err != nil {
		t.Error(err)
	}

	if !result.Package.ContainsOnlyChanges {
		t.Error("ContainsOnlyChanges (xml attribute 'СодержитТолькоИзменения') failed to parse (should be true)")
	}

	if len(result.Package.Offers) != 3 {
		t.Fatalf("Offers failed to parse, there should be 3 offers but found %d", len(result.Package.Offers))
	}

	if len(result.Package.PriceTypes) != 1 {
		t.Fatalf("Failed to parse PriceTypes found %d, expected 1", len(result.Package.PriceTypes))
	}

}

var (
	offersData = []byte(`
	<?xml version="1.0" encoding="UTF-8"?>
	<КоммерческаяИнформация xmlns="urn:1C.ru:commerceml_210" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" ВерсияСхемы="2.08" ДатаФормирования="2025-06-24T21:54:40">
		<ПакетПредложений СодержитТолькоИзменения="true">
			<ТипыЦен>
				<ТипЦены>
					<Ид>84d7160b-7c8d-11ed-bbdb-e0b9a548d6d8</Ид>
					<Наименование>Розничная цена</Наименование>
					<Валюта>руб</Валюта>
					<Налог>
						<Наименование>НДС</Наименование>
						<УчтеноВСумме>true</УчтеноВСумме>
					</Налог>
				</ТипЦены>
			</ТипыЦен>
			<Предложения>
				<Предложение>
					<Ид>7c63329c-7b37-11ee-802b-e0b9a548d6d8</Ид>
					<Наименование>Авт. выкл. 1Р 6А х-ка C ВА-103 6кА DEKraft</Наименование>
					<БазоваяЕдиница Код="796" НаименованиеПолное="штука" МеждународноеСокращение="PCE"/>
					<Артикул>12269DEK</Артикул>
					<Цены>
						<Цена>
							<Представление>296,4 руб. за шт</Представление>
							<ИдТипаЦены>84d7160b-7c8d-11ed-bbdb-e0b9a548d6d8</ИдТипаЦены>
							<ЦенаЗаЕдиницу>296.4</ЦенаЗаЕдиницу>
							<Валюта>руб</Валюта>
							<Единица>шт</Единица>
							<Коэффициент>1</Коэффициент>
						</Цена>
					</Цены>
					<Количество>220</Количество>
				</Предложение>
				<Предложение>
					<Ид>a09d2dd8-7fa0-11ed-a9dc-e0b9a548d6d8</Ид>
					<Наименование>Шкаф распределительный</Наименование>
					<БазоваяЕдиница Код="796" НаименованиеПолное="штука" МеждународноеСокращение="PCE"/>
					<Артикул>MPS 220.80.80</Артикул>
					<Цены>
						<Цена>
							<Представление>96 287 руб. за шт</Представление>
							<ИдТипаЦены>84d7160b-7c8d-11ed-bbdb-e0b9a548d6d8</ИдТипаЦены>
							<ЦенаЗаЕдиницу>96287</ЦенаЗаЕдиницу>
							<Валюта>руб</Валюта>
							<Единица>шт</Единица>
							<Коэффициент>1</Коэффициент>
						</Цена>
					</Цены>
					<Количество>0</Количество>
				</Предложение>
				<Предложение>
					<Ид>b2d23a47-7cdb-11ef-ad08-000c292ac68f</Ид>
					<Наименование>КНЗ 2,5-2-PE  Заземляющая клемма  2,5мм</Наименование>
					<БазоваяЕдиница Код="796" НаименованиеПолное="штука" МеждународноеСокращение="PCE"/>
					<Артикул>10000012</Артикул>
					<Цены>
						<Цена>
							<Представление>170,04 руб. за шт</Представление>
							<ИдТипаЦены>84d7160b-7c8d-11ed-bbdb-e0b9a548d6d8</ИдТипаЦены>
							<ЦенаЗаЕдиницу>170.04</ЦенаЗаЕдиницу>
							<Валюта>руб</Валюта>
							<Единица>шт</Единица>
							<Коэффициент>1</Коэффициент>
						</Цена>
					</Цены>
					<Количество>120</Количество>
				</Предложение>
			</Предложения>
		</ПакетПредложений>
	</КоммерческаяИнформация>`)

	importsData = []byte(`
	<КоммерческаяИнформация>
	<Классификатор>
	<Категории>
		<Категория>
			<Ид>b0627bbe-ceb3-11eb-a53c-f14c52115a4f</Ид>
			<Наименование>&lt;Без категории&gt;</Наименование>
			<Свойства>
				<Ид>196f8eb6-7f84-11ed-a9dc-e0b9a548d6d8</Ид>
				<Ид>31036664-469c-11f0-ad0f-000c292ac68f</Ид>
				<Ид>12d51e24-469d-11f0-ad0f-000c292ac68f</Ид>
				<Ид>41838cbb-469d-11f0-ad0f-000c292ac68f</Ид>
				<Ид>b7f236e2-469e-11f0-ad0f-000c292ac68f</Ид>
				<Ид>133a5a26-469f-11f0-ad0f-000c292ac68f</Ид>
				<Ид>74d6567c-469f-11f0-ad0f-000c292ac68f</Ид>
				<Ид>b86c05fe-46a0-11f0-ad0f-000c292ac68f</Ид>
				<Ид>e7865f1a-46a0-11f0-ad0f-000c292ac68f</Ид>
				<Ид>0abc3925-46a1-11f0-ad0f-000c292ac68f</Ид>
				<Ид>7338a8d5-46a1-11f0-ad0f-000c292ac68f</Ид>
				<Ид>9a965d49-46a1-11f0-ad0f-000c292ac68f</Ид>
				<Ид>d9772f10-46a1-11f0-ad0f-000c292ac68f</Ид>
				<Ид>2cc436e9-46a3-11f0-ad0f-000c292ac68f</Ид>
				<Ид>ac448f94-46a3-11f0-ad0f-000c292ac68f</Ид>
			</Свойства>
		</Категория>
	</Категории>
	</Классификатор>
		<Каталог СодержитТолькоИзменения="true">
			<Товары>
				<Товар>
					<Ид>RandomId</Ид>
					<Категория>b0627bbe-ceb3-11eb-a53c-f14c52115a4f</Категория>
					<Группы>
						<Ид>GroupId</Ид>
					</Группы>
					<ЗначенияСвойств>
						<ЗначенияСвойства>
							<Ид>PropValId</Ид>
							<Значение>PropValVal</Значение>
						</ЗначенияСвойства>
					</ЗначенияСвойств>
				</Товар>
			</Товары>
		</Каталог>
	</КоммерческаяИнформация>`)
)
