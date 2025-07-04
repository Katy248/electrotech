package models

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyHandbookValues = errors.New("handbook values must not be empty")
)

type ParameterType string

const (
	ParameterTypeNumber ParameterType = "number"
	ParameterTypeList   ParameterType = "list"
)

type Parameter struct {
	Name     string        `json:"name"`
	Type     ParameterType `json:"type"`
	Values   []string      `json:"values,omitzero"`
	MinValue float64       `json:"minValue,omitzero"`
	MaxValue float64       `json:"maxValue,omitzero"`
}

func NewListParameter(name string, values []string) (*Parameter, error) {
	if len(values) == 0 {
		return nil, ErrEmptyHandbookValues
	}
	return &Parameter{
		Type:   ParameterTypeList,
		Values: values,
	}, nil
}

func NewNumberParameter(name string, min, max float64) (*Parameter, error) {
	if min > max {
		return nil, fmt.Errorf("min value must be less than or equal to max value (min = %f, max = %f)", min, max)
	}
	return &Parameter{
		Type:     ParameterTypeNumber,
		Name:     name,
		MinValue: min,
		MaxValue: max,
	}, nil
}
