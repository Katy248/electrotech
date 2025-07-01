package parser

import (
	"encoding/xml"
	"errors"
	"os"
)

type Parser struct {
	dir string
}

var (
	ErrOffersFileNotFound  = errors.New("offers file not found")
	ErrImportsFileNotFound = errors.New("imports file not found")
)

func NewParser(directory string) (*Parser, error) {

	if !fileExists(getOffersFilepath(directory)) {
		return nil, ErrOffersFileNotFound
	}
	if !fileExists(getImportsFilepath(directory)) {
		return nil, ErrImportsFileNotFound
	}
	return &Parser{dir: directory}, nil
}

func getOffersFilepath(dir string) string {
	return dir + "/offers.xml"
}

func getImportsFilepath(dir string) string {
	return dir + "/import.xml"
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (p *Parser) Parse() error {
	return nil
}

func parseImportsData(data []byte) (*ImportsModel, error) {
	var model ImportsModel
	err := xml.Unmarshal(data, &model)

	return &model, err
}

func parseOffersData(data []byte) (*OffersModel, error) {
	var model OffersModel
	err := xml.Unmarshal(data, &model)

	return &model, err
}
