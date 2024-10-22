package forms

import (
	"net/url"

	"github.com/PsionicAlch/psionicalch-home/internal/validators"
)

type SignUpForm struct {
	Email           string
	Password        string
	ConfirmPassword string
	RememberMe      bool
	Errors          FormErrors
}

func (form *SignUpForm) GetErrors() FormErrors {
	return form.Errors
}

func (form *SignUpForm) SetErrors(errs FormErrors) {
	form.Errors = errs
}

func NewSignupFormErrors() FormErrors {
	return map[string][]string{
		"email":            {},
		"password":         {},
		"confirm_password": {},
	}
}

func NewSignupForm() *SignUpForm {
	email := ""
	password := ""
	confirmPassword := ""
	rememberMe := false
	formErrors := NewSignupFormErrors()

	return &SignUpForm{
		Email:           email,
		Password:        password,
		ConfirmPassword: confirmPassword,
		RememberMe:      rememberMe,
		Errors:          formErrors,
	}
}

func CreateSignUpForm(form url.Values) *SignUpForm {
	email := form.Get("email")
	password := form.Get("password")
	confirmPassword := form.Get("confirm-password")
	rememberMe := form.Has("remember-me")
	formErrors := NewSignupFormErrors()

	return &SignUpForm{
		Email:           email,
		Password:        password,
		ConfirmPassword: confirmPassword,
		RememberMe:      rememberMe,
		Errors:          formErrors,
	}
}

func (form *SignUpForm) Validate() bool {
	valid := true

	if err := validators.ValidateEmail(form.Email); err != nil {
		AppendErrors(err, "email", form)
		valid = false
	}

	if err := validators.ValidatePassword(form.Password, 8); err != nil {
		AppendErrors(err, "password", form)
		valid = false
	}

	if err := validators.ValidatePasswordsMatch(form.Password, form.ConfirmPassword); err != nil {
		AppendErrors(err, "confirm_password", form)
		valid = false
	}

	return valid
}

func (form *SignUpForm) ValidateWithoutEmpty() bool {
	valid := true

	if err := validators.ValidateEmailWithoutEmpty(form.Email); err != nil {
		AppendErrors(err, "email", form)
		valid = false
	}

	if err := validators.ValidatePasswordWithoutEmpty(form.Password, 8); err != nil {
		AppendErrors(err, "password", form)
		valid = false
	}

	if err := validators.ValidatePasswordsMatch(form.Password, form.ConfirmPassword); err != nil {
		AppendErrors(err, "confirm_password", form)
		valid = false
	}

	return valid
}

func (form *SignUpForm) IsValid() bool {
	return len(form.Errors) == 0
}
