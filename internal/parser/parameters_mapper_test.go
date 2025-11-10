package parser

import (
	"os"
	"testing"
)

func TestMapParameters(t *testing.T) {

	dir, _ := os.Getwd()
	t.Logf("Current directory: %s", dir)

	p, err := NewParser("../../example")
	if err != nil {
		t.Fatal(err)
	}

	if err := p.parse(); err != nil {
		t.Fatalf("Failed parse catalog: %s", err)
	}

	parameters, err := p.mapParameters()
	if err != nil {
		t.Fatalf("Failed map parameters: %s", err)
	}

	t.Logf("Parameters: %v", parameters)
}

func TestParseFloat(t *testing.T) {

	testCase := func(value string, expected float64) {
		f := parseFloat(value)
		if f != expected {
			t.Errorf("Failed parse '%s', got %f, expected %f", value, f, expected)
		}
	}

	testCase("1.23", 1.23)
	testCase("1,23", 1.23)
	testCase("", 0)
	testCase("0", 0)
	testCase("wsdgfvjsdnijsnvjkdnkj", 0)
	testCase("0,3", 0.3)
}
