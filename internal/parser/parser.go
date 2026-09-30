package parser

import (
	"electrotech/internal/models"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path"

	"charm.land/log/v2"
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

	return &Parser{dir: directory}, nil //nolint:exhaustruct_v5
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

func (p *Parser) GetProducts() ([]models.Product, error) {
	if err := p.parse(); err != nil {
		return nil, fmt.Errorf("failed parse xml data: %w", err)
	}

	products, err := mapProducts(p.offers, p.imports)
	if err != nil {
		return nil, fmt.Errorf("failed map xml data: %w", err)
	}

	return products, nil
}
func (p *Parser) parse() error {
	if p.imports == nil {
		imp, err := p.parseImports()
		if err != nil {
			return fmt.Errorf("failed parse imports: %w", err)
		}

		p.imports = imp
	}

	if p.offers == nil {
		off, err := p.parseOffers()
		if err != nil {
			return fmt.Errorf("failed parse offers: %w", err)
		}

		p.offers = off
	}

	return nil
}

const RootDir = "."

func getDataFromFile(filepath string) ([]byte, error) {
	filepath = path.Clean(filepath)

	root, err := os.OpenRoot(RootDir)
	if err != nil {
		return nil, fmt.Errorf("open root %q: %w", RootDir, err)
	}
	defer func() {
		if err := root.Close(); err != nil {
			log.Error("Failed close root dir", "dir", RootDir, "error", err)
		}
	}()

	data, err := root.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("read file %q: %w", filepath, err)
	}

	return data, nil
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
	if err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	return &model, nil
}

func parseOffersData(data []byte) (*offersModel, error) {
	var model offersModel

	err := xml.Unmarshal(data, &model)
	if err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	return &model, nil
}
