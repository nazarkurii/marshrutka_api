package credentianlsutils

import (
	"errors"

	"github.com/nyaruka/phonenumbers"
)

func FomratPhoneNumber(phoneNumber string) (string, error) {
	pn, err := phonenumbers.Parse(phoneNumber, "UA")
	if err != nil {
		return "", err
	}

	if !phonenumbers.IsValidNumber(pn) {
		return "", errors.New("invalid phone number")
	}

	return phonenumbers.Format(pn, phonenumbers.E164), nil
}
