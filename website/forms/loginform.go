package forms

import (
	"net/url"

	"github.com/PsionicAlch/psionicalch-home/internal/validators"
)

type LoginForm struct {
	Email      string
	Password   string
	RememberMe bool
	Errors     FormErrors
}

func (form *LoginForm) GetErrors() FormErrors {
	return form.Errors
}

func (form *LoginForm) SetErrors(errs FormErrors) {
	form.Errors = errs
}

func NewLoginFormErrors() FormErrors {
	return map[string][]string{
		"email":    {},
		"password": {},
	}
}

func NewLoginForm() *LoginForm {
	loginForm := new(LoginForm)
	loginForm.RememberMe = false
	loginForm.Errors = NewLoginFormErrors()

	return loginForm
}

func CreateLoginForm(form url.Values) *LoginForm {
	loginForm := NewLoginForm()

	loginForm.Email = form.Get("email")
	loginForm.Password = form.Get("password")
	loginForm.RememberMe = form.Has("remember-me")

	return loginForm
}

func (form *LoginForm) Validate() bool {
	valid := true

	if err := validators.ValidateEmail(form.Email); err != nil {
		AppendErrors(err, "email", form)
		valid = false
	}

	if err := validators.ValidateNotEmpty(form.Password); err != nil {
		AppendErrors(err, "password", form)
		valid = false
	}

	return valid
}

func (form *LoginForm) ValidateWithoutEmpty() bool {
	valid := true

	if err := validators.ValidateEmailWithoutEmpty(form.Email); err != nil {
		AppendErrors(err, "email", form)
		valid = false
	}

	return valid
}

func (form *LoginForm) IsValid() bool {
	return len(form.Errors) == 0
}

func (form *LoginForm) AddError(field, msg string) {
	form.Errors[field] = append(form.Errors[field], msg)
}
