package catalog

import (
	"electrotech/internal/parser"
	"fmt"

	"github.com/spf13/viper"
)

var (
	ErrNotImplemented      = fmt.Errorf("this function is not implemented")
	ErrDataDirNotSpecified = fmt.Errorf("environment variable DATA_DIR isn't specified")
)

type Page int

type Repo struct {
	parser *parser.Parser
}

func New() (*Repo, error) {
	dataDir := viper.GetString("data-dir")
	if dataDir == "" {
		return nil, ErrDataDirNotSpecified
	}

	p, err := parser.NewParser(dataDir)
	if err != nil {
		return nil, err
	}
	return &Repo{parser: p}, nil
}
