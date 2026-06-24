package catalog

import (
	"electrotech/internal/parser"
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

var (
	ErrNotImplemented      = errors.New("this function is not implemented")
	ErrDataDirNotSpecified = errors.New("data-dir parameter isn't specified")
)

// Page is a page number for pagination.
// From 0 to infinity.
type Page int

type Repo struct {
	parser *parser.Parser
}

func New() (*Repo, error) {
	viper.SetDefault("data-dir", "/data")

	dataDir := viper.GetString("data-dir")
	if dataDir == "" {
		return nil, ErrDataDirNotSpecified
	}

	p, err := parser.NewParser(dataDir)
	if err != nil {
		return nil, fmt.Errorf("new parser: %w", err)
	}

	return &Repo{parser: p}, nil
}
