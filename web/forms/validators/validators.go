package validators

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	emailverifier "github.com/AfterShip/email-verifier"
	goaway "github.com/TwiN/go-away"
)

type ValidationFunc func(data string, values url.Values) error

func ChainValidators(validatorFuncs ...ValidationFunc) ValidationFunc {
	return func(data string, values url.Values) error {
		for _, validator := range validatorFuncs {
			err := validator(data, values)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func Empty(data string, values url.Values) error {
	return nil
}

func NotEmpty(data string, values url.Values) error {
	if data == "" {
		return errors.New("cannot be empty")
	}

	return nil
}

func MaxLength(length uint) ValidationFunc {
	return func(data string, values url.Values) error {
		if data == "" {
			if len(data) > int(length) {
				return fmt.Errorf("can't be more than %d characters long", length)
			}
		}

		return nil
	}
}

func MinLength(length uint) ValidationFunc {
	return func(data string, values url.Values) error {
		if data != "" {
			if len(data) < int(length) {
				return fmt.Errorf("can't be less than %d characters long", length)
			}
		}

		return nil
	}
}

func Profanities(data string, values url.Values) error {
	if data != "" {
		if goaway.IsProfane(data) {
			return errors.New("can't contain any profanities")
		}
	}

	return nil
}

func NoSpecialCharacters(data string, values url.Values) error {
	if data != "" {
		if ContainsSpecialCharacters(data) {
			return errors.New("can't contain special characters")
		}
	}

	return nil
}

func NoNumbers(data string, values url.Values) error {
	if data != "" {
		if ContainsNumbers(data) {
			return errors.New("can't contain numbers")
		}
	}

	return nil
}

var emailVerifier = emailverifier.NewVerifier()

func IsEmail(data string, values url.Values) error {
	if data != "" {
		errInvalidEmail := errors.New("email isn't valid")

		result, err := emailVerifier.Verify(data)
		if err != nil {
			return errInvalidEmail
		}

		if !result.Syntax.Valid {
			return errInvalidEmail
		}
	}

	return nil
}

func IsNotDisposableEmail(data string, values url.Values) error {
	if data != "" {
		errInvalidEmail := errors.New("email isn't valid")
		errDisposableEmail := errors.New("email can't be disposable")

		result, err := emailVerifier.Verify(data)
		if err != nil {
			return errInvalidEmail
		}

		if !result.Syntax.Valid {
			return errInvalidEmail
		}

		if result.Disposable {
			return errDisposableEmail
		}
	}

	return nil
}

func UppercaseCharacters(data string, values url.Values) error {
	if data != "" {
		if !ContainsUppercaseCharacters(data) {
			return errors.New("needs to have at least one uppercase character")
		}
	}

	return nil
}

func LowercaseCharacters(data string, values url.Values) error {
	if data != "" {
		if !ContainsLowercaseCharacters(data) {
			return errors.New("needs to have at least one lowercase character")
		}
	}

	return nil
}

func NumberCharacters(data string, values url.Values) error {
	if data != "" {
		if !ContainsNumbers(data) {
			return errors.New("needs to have at least one number")
		}
	}

	return nil
}

func MatchesField(field, fieldName string) ValidationFunc {
	return func(data string, values url.Values) error {
		if data != values.Get(field) {
			return fmt.Errorf("not the same as %s", fieldName)
		}

		return nil
	}
}

func Integer(data string, values url.Values) error {
	if _, err := strconv.Atoi(data); err != nil {
		return errors.New("needs to be a number")
	}

	return nil
}

func Max(max int) ValidationFunc {
	return func(data string, values url.Values) error {
		num, err := strconv.Atoi(data)
		if err != nil {
			return errors.New("needs to be a number")
		}

		if num > max {
			return fmt.Errorf("cannot be bigger than %d", max)
		}

		return nil
	}
}

func Min(min int) ValidationFunc {
	return func(data string, values url.Values) error {
		num, err := strconv.Atoi(data)
		if err != nil {
			return errors.New("needs to be a number")
		}

		if num < min {
			return fmt.Errorf("cannot be smaller than %d", min)
		}

		return nil
	}
}
