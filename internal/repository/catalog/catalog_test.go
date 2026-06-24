package catalog_test

import (
	"electrotech/internal/repository/catalog"
	"errors"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestNewCatalogWithoutEnv(t *testing.T) {
	t.Parallel()
	viper.Set("data-dir", "")

	_, err := catalog.New()
	if !errors.Is(err, catalog.ErrDataDirNotSpecified) {
		t.Errorf("Expected '%s' error but there is '%s'", catalog.ErrDataDirNotSpecified, err)
	}
}
func TestNewCatalogBadDir(t *testing.T) {
	t.Setenv("DATA_DIR", "./not-exist")

	_, err := catalog.New()
	if err == nil {
		t.Error("There is not error, but shuld be, cause directory not exist")
	}
}

func TestNewCatalog(t *testing.T) {
	t.Parallel()

	currentDir, _ := os.Getwd()
	t.Logf("Current dir: %s", currentDir)
	viper.Set("data-dir", "../../../example")

	_, err := catalog.New()
	if err != nil {
		t.Errorf("Failed create repository: %s", err)
	}
}
