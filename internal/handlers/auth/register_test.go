package auth

import "testing"

func TestFormatPhoneNumber(t *testing.T) {
	testCases := []struct {
		input         string
		expected      string
		shouldBeError bool
	}{
		{"+79991234567", "89991234567", false},
		{"79991234567", "89991234567", true},
		{"32 9991234567", "89991234567", true},
		{
			input:    "+7 (999) 123-45-67",
			expected: "89991234567"},
		{
			input:         "",
			shouldBeError: true,
		},
	}
	for _, test := range testCases {
		actual, err := FormatPhoneNumber(test.input)
		if err != nil {
			if !test.shouldBeError {
				t.Errorf("FormatPhoneNumber(%q) returned error %q, but should not", test.input, err)
			}
		} else {
			if actual != test.expected {
				t.Errorf("FormatPhoneNumber(%q) = %q, expected %q", test.input, actual, test.expected)
			}
		}
	}
}
