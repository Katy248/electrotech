package catalog

import (
	"os"
	"testing"
)

func TestNewCatalogWithoutEnv(t *testing.T) {
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

func TestNewCatalog(t *testing.T) {
	currentDir, _ := os.Getwd()
	t.Logf("Current dir: %s", currentDir)
	os.Setenv("DATA_DIR", "../../../example")
	_, err := New()
	if err != nil {
		t.Errorf("Failed create repository: %s", err)
	}
}
