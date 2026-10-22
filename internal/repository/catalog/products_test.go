package catalog_test

import (
	"context"
	"electrotech/internal/repository/catalog"
	"testing"
)

func TestGetPages(t *testing.T) {
	t.Parallel()

	if result := catalog.GetPages(10, 10); result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}

	if result := catalog.GetPages(21, 20); result != 2 {
		t.Errorf("Expected 2, got %d", result)
	}
}
func TestGetProduct(t *testing.T) {
	t.Parallel()

	c, err := catalog.New(NewConfig(t, "../../../example"))
	if err != nil {
		t.Errorf("Failed create repository: %s", err)
	}

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	productID := "a09d2dd8-7fa0-11ed-a9dc-e0b9a548d6d8"

	_, err = c.GetProduct(ctx, productID)
	if err == nil {
		t.Errorf("There must be error, but wasn't")
	}
}
