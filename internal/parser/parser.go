package parser

import (
	"electrotech/internal/models"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
)

type Parser struct {
	dir string

	offers  *offersModel
	imports *importsModel
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

func (p *Parser) parse() error {
	if p.imports == nil {
		imp, err := p.parseImports()
		if err != nil {
			return err
		}
		p.imports = imp
	}

	if p.offers == nil {
		off, err := p.parseOffers()
		if err != nil {
			return err
		}
		p.offers = off
	}
	return nil
}

func (p *Parser) GetProducts() ([]models.Product, error) {
	if err := p.parse(); err != nil {
		return nil, fmt.Errorf("failed parse data: %s", err)
	}

	products, err := mapProducts(p.offers, p.imports)
	if err != nil {
		return nil, fmt.Errorf("failed map xml data: %s", err)
	}

	return products, nil
}

func getDataFromFile(filepath string) ([]byte, error) {
	data, err := os.ReadFile(filepath)
	return data, err
}

func (p *Parser) parseImports() (*importsModel, error) {
	data, err := getDataFromFile(getImportsFilepath(p.dir))
	if err != nil {
		return nil, err
	}

	return parseImportsData(data)
}

func (p *Parser) parseOffers() (*offersModel, error) {
	data, err := getDataFromFile(getOffersFilepath(p.dir))
	if err != nil {
		return nil, err
	}

	return parseOffersData(data)
}

func parseImportsData(data []byte) (*importsModel, error) {
	var model importsModel
	err := xml.Unmarshal(data, &model)

	return &model, err
}

func parseOffersData(data []byte) (*offersModel, error) {
	var model offersModel
	err := xml.Unmarshal(data, &model)

	return &model, err
}
