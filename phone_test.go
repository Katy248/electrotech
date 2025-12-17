package electrotech

import "testing"

func TestFormatPhoneNumber(t *testing.T) {
	testCases := []struct {
		input         string
		expected      string
		shouldBeError bool
	}{
		{
			input:    "+79991234567",
			expected: "+79991234567",
		},
		{
			input:         "79991234567",
			shouldBeError: true,
		},
		{
			input:         "1232 9991234567",
			shouldBeError: true,
		},
		{
			input:    "+7 (999) 123-45-67",
			expected: "+79991234567",
		},
		{
			input:    "8 (999) 123-45-67",
			expected: "+79991234567",
		},
		{
			input:         "",
			shouldBeError: true,
		},
		{
			input:         " 8 8   8   8  888   ",
			shouldBeError: true,
		},
		{
			input:         "some simple chars, ---",
			shouldBeError: true,
		},
		{
			input:         "+7 958 803 92-95c",
			shouldBeError: true,
		},
	}
	for _, test := range testCases {
		actual, err := FormatPhoneNumber(test.input)
		if test.shouldBeError && err == nil {
			t.Errorf("FormatPhoneNumber(%q) returned no error, but should have", test.input)
			break
		}

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
