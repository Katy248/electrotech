package catalog_test

import (
	"electrotech/internal/repository/catalog"
	"errors"
	"os"
	"testing"
)

func NewConfig(t *testing.T, dataDir string) *catalog.Config {
	t.Helper()

	return &catalog.Config{
		DataDir: dataDir,
	}
}

func TestNewCatalogWithoutEnv(t *testing.T) {
	t.Parallel()

	_, err := catalog.New(NewConfig(t, ""))
	if !errors.Is(err, catalog.ErrDataDirNotSpecified) {
		t.Errorf("Expected '%s' error but there is '%s'", catalog.ErrDataDirNotSpecified, err)
	}
}
func TestNewCatalogBadDir(t *testing.T) {
	t.Parallel()

	_, err := catalog.New(NewConfig(t, "./not-exist"))
	if err == nil {
		t.Error("There is not error, but shuld be, cause directory not exist")
	}
}

func TestNewCatalog(t *testing.T) {
	t.Parallel()

	currentDir, _ := os.Getwd()
	t.Logf("Current dir: %s", currentDir)

	_, err := catalog.New(NewConfig(t, "../../../example"))
	if err != nil {
		t.Errorf("Failed create repository: %s", err)
	}
}
