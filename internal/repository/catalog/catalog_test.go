package catalog

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestNewCatalogWithoutEnv(t *testing.T) {
	viper.Set("data-dir", "")
	_, err := New()
	if err != ErrDataDirNotSpecified {
		t.Errorf("Expected '%s' error but there is '%s'", ErrDataDirNotSpecified, err)
	}
}
func TestNewCatalogBadDir(t *testing.T) {
	os.Setenv("DATA_DIR", "./not-exist")
	_, err := New()
	if err == nil {
		t.Error("There is not error, but shuld be, cause directory not exist")
	}
}

func testGetNewCatalog(t *testing.T) *Repo {
	currentDir, _ := os.Getwd()
	t.Logf("Current dir: %s", currentDir)
	os.Setenv("DATA_DIR", "../../../example")
	catalog, err := New()
	if err != nil {
		t.Fatalf("Failed create repository: %s", err)
	}
	return catalog
}

func TestNewCatalog(t *testing.T) {
	currentDir, _ := os.Getwd()
	t.Logf("Current dir: %s", currentDir)
	viper.Set("data-dir", "../../../example")
	_, err := New()
	if err != nil {
		t.Errorf("Failed create repository: %s", err)
	}
}
