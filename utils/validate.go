package utils

import (
	"net/mail"
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

var regPhone = regexp.MustCompile(`\+\d{1,4}\d{1,10}\b`)

func ValidatePhone(phone string) (string, error) {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, ".", "")

	match := regPhone.MatchString(phone)
	if !match {
		return "", ErrorInvalidPhone
	}

	return phone, nil
}
