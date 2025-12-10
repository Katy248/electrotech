package catalog

import "testing"

func TestGetPages(t *testing.T) {
	if result := getPages(10, 10); result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}
	if result := getPages(21, 20); result != 2 {
		t.Errorf("Expected 2, got %d", result)
	}
}
