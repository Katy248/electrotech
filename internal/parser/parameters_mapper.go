package parser

import (
	"electrotech/internal/models"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
)

func (p *Parser) GetParameters() ([]models.Parameter, error) {
	if err := p.parse(); err != nil {
		return nil, fmt.Errorf("failed parse catalog: %s", err)
	}
	return p.mapParameters()
}

func (p *Parser) mapParameters() ([]models.Parameter, error) {
	var result []models.Parameter
	for _, sourceParam := range p.imports.Classifier.Properties {
		switch sourceParam.Type {
		case propertyTypeHandbook:
			values := variantsToValues(sourceParam.Variants)
			log.Debug("Handbook parameter", "param", sourceParam.Name, "values", values)
			handbook, err := models.NewListParameter(sourceParam.Name, values)
			if err != nil {
				return result, err
			}
			result = append(result, *handbook)
		case propertyTypeNumber:
			min, max, err := getDiapason(sourceParam.Id, p.imports.Catalog.Products)
			if err != nil {
				return result, err
			}
			number, err := models.NewNumberParameter(sourceParam.Name, min, max)
			if err != nil {
				return result, err
			}
			result = append(result, *number)
		case propertyTypeString:
			values := getStringParameterEntries(sourceParam.Id, p.imports.Catalog.Products)
			p, err := models.NewListParameter(sourceParam.Name, values)
			if err != nil {
				return result, err
			}
			result = append(result, *p)
		default:
			return result, fmt.Errorf("unknown parameter type '%s' for parameter '%s'", sourceParam.Type, sourceParam.Name)

		}
	}
	return result, nil
}

func variantsToValues(variants []handbookVariant) []string {
	var values []string
	for _, variant := range variants {
		values = append(values, variant.Value)
	}
	return values
}

func getStringParameterEntries(propId string, products []product) []string {
	var result []string
	for _, p := range products {
		for _, param := range p.PropertyValues {
			if param.Id != propId {
				continue
			}
			if slices.Contains(result, param.Value) {
				continue
			}
			result = append(result, param.Value)
		}
	}
	return result

}

// TODO: TEST TEST TEST
// There should be tests
// Please add tests ASAP! Or call Katy248 for this
func getDiapason(propId string, products []product) (float64, float64, error) {

	min := float64(math.MaxInt64)
	max := float64(math.MinInt64)

	for _, p := range products {
		for _, param := range p.PropertyValues {
			if param.Id != propId {
				continue
			}
			intValue := parseFloat(param.Value)

			if intValue < min {
				min = intValue
			}
			if intValue > max {
				max = intValue
			}
		}
	}
	return min, max, nil
}

func parseFloat(value string) float64 {
	if value == "" {
		return 0
	}
	value = strings.ReplaceAll(value, ",", ".")
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return result
}
