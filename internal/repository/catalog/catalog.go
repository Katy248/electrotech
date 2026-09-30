package catalog

import (
	"electrotech/internal/parser"
	"errors"
	"fmt"
)

var (
	ErrNotImplemented      = errors.New("this function is not implemented")
	ErrDataDirNotSpecified = errors.New("data directory isn't specified")
)

type Config struct {
	DataDir string `mapstructure:"data-dir"`
}

// Page is a page number for pagination.
// From 0 to infinity.
type Page int

type Repo struct {
	parser *parser.Parser
}

func New(config *Config) (*Repo, error) {
	if config.DataDir == "" {
		return nil, ErrDataDirNotSpecified
	}

	p, err := parser.NewParser(config.DataDir)
	if err != nil {
		return nil, fmt.Errorf("new parser: %w", err)
	}

	return &Repo{parser: p}, nil
}
