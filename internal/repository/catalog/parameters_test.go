package catalog

import "testing"

func TestGetParameters(t *testing.T) {

	catalog := testGetNewCatalog(t)
	params, err := catalog.GetParameters()
	if err != nil {
		t.Fatal(err)
	}
	if len(params) == 0 {
		t.Error("No parameters found")
	}

	t.Logf("Parameters: %v", params)
}
