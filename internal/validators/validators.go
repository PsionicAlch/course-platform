package validators

import (
	"fmt"
	"net/mail"
	"regexp"
	"unicode/utf8"
)

type ValidationError struct {
	Errors []string
}

func CreateValidationError(errorMsgs []string) *ValidationError {
	return &ValidationError{
		Errors: errorMsgs,
	}
}

func (v *ValidationError) Error() string {
	var errorMsgs string

	for _, err := range v.Errors {
		errorMsgs = fmt.Sprintf("%s. %s", errorMsgs, err)
	}

	return errorMsgs
}

// EmptyString checks if a string is empty.
func EmptyString(s string) bool {
	return len(s) == 0 && s == ""
}

func ValidateEmailWithoutEmpty(email string) error {
	if !EmptyString(email) {
		em, err := mail.ParseAddress(email)
		if !(err == nil && em.Address == email) {
			return CreateValidationError([]string{"your email is invalid"})
		}
	}

	return nil
}

// ValidateEmail checks to make sure that an email address is valid.
func ValidateEmail(email string) error {
	errMsgs := make([]string, 0, 1)

	if EmptyString(email) {
		errMsgs = append(errMsgs, "you cannot have an empty email")
		return CreateValidationError(errMsgs)
	}

	return ValidateEmailWithoutEmpty(email)
}

func ValidatePasswordWithoutEmpty(password string, length int) error {
	if !EmptyString(password) {
		errMsgs := make([]string, 0, 5)

		if utf8.RuneCountInString(password) < length {
			errMsgs = append(errMsgs, fmt.Sprintf("your password needs to be at least %d characters long", length))
		}

		if match, err := regexp.MatchString("[A-Z]", password); err != nil || !match {
			errMsgs = append(errMsgs, "your password needs to contain at least one capital letter")
		}

		if match, err := regexp.MatchString("[a-z]", password); err != nil || !match {
			errMsgs = append(errMsgs, "your password needs to contain at least one lowercase letter")
		}

		if match, err := regexp.MatchString("[0-9]", password); err != nil || !match {
			errMsgs = append(errMsgs, "your password needs to contain at least one number")
		}

		if match, err := regexp.MatchString(`[!@#$%^&*(),.?":{}|<>]`, password); err != nil || !match {
			errMsgs = append(errMsgs, "your password needs to contain at least one special character")
		}

		if len(errMsgs) > 0 {
			return CreateValidationError(errMsgs)
		}
	}

	return nil
}

// ValidatePassword checks to make sure that the password is valid.
func ValidatePassword(password string, length int) error {
	errMsgs := make([]string, 0, 5)

	if EmptyString(password) {
		errMsgs = append(errMsgs, "password cannot be empty")
		return CreateValidationError(errMsgs)
	}

	if err := ValidatePasswordWithoutEmpty(password, length); err != nil {
		validationErr, ok := err.(*ValidationError)
		if ok {
			errMsgs = append(errMsgs, validationErr.Errors...)
		}
	}

	if len(errMsgs) > 0 {
		return CreateValidationError(errMsgs)
	}

	return nil
}

// ValidatePasswordsMatch checks to make sure that both passwords are the same.
func ValidatePasswordsMatch(password1, password2 string) error {
	if password1 != password2 {
		return CreateValidationError([]string{"the passwords you entered don't match"})
	}

	return nil
}

func ValidateNotEmpty(value string) error {
	if value == "" || len(value) <= 0 {
		return CreateValidationError([]string{"cannot be empty"})
	}

	return nil
}
