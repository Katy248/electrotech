package electrotech

import (
	"errors"
	"fmt"
	"strings"
)

var ErrPhoneNumberInvalid = errors.New("invalid phone number")
var ErrPhoneNumberEmpty = fmt.Errorf("%w: phone number is empty", ErrPhoneNumberInvalid)

const PhoneNumberMinLength = 11
const PhoneNumberFormattedLength = 12

func FormatPhoneNumber(phone string) (string, error) {
	if phone == "" {
		return "", ErrPhoneNumberEmpty
	}

	if len(phone) < PhoneNumberMinLength {
		return "", fmt.Errorf(
			"%w: probably invalid phone number, length is less than 11 (%d)",
			ErrPhoneNumberInvalid,
			len(phone),
		)
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
			return phone, fmt.Errorf("%w: invalid character %q at index %d", ErrPhoneNumberInvalid, ch, index)
		}
	}

	if len(phone) != PhoneNumberFormattedLength {
		return phone, fmt.Errorf("%w: invalid phone number length %d", ErrPhoneNumberInvalid, len(phone))
	}

	return phone, nil
}
