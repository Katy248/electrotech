package electrotech

import (
	"fmt"
	"strings"
)

func FormatPhoneNumber(phone string) (string, error) {
	if phone == "" {
		return "", fmt.Errorf("phone number is empty")

	}
	if len(phone) < 11 {
		return "", fmt.Errorf("probably invalid phone number, length is less than 11 (%d)", len(phone))
	}
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	phone = strings.ReplaceAll(phone, " ", "")

	if phone[0] == '8' {
		phone = "+7" + phone[1:]
	}

	for index, ch := range phone {
		if ch != '+' && ch < '0' || ch > '9' {
			return phone, fmt.Errorf("invalid character %q at index %d", ch, index)
		}
	}

	if len(phone) != 12 {
		return phone, fmt.Errorf("invalid phone number length %d", len(phone))
	}
	return phone, nil
}
