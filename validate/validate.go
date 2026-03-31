package validate

import (
	"errors"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Email validates that the given string is a valid email address format.
func Email(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{6,14}$`)

// Phone validates that the given string is a valid international phone number (E.164 format).
func Phone(phone string) error {
	if !phoneRegex.MatchString(phone) {
		return errors.New("invalid phone number format")
	}
	return nil
}
