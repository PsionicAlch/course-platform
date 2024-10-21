package forms

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/validators"
)

type SignUpForm struct {
	email           string
	password        string
	confirmPassword string
	rememberMe      bool
	formErrors      map[string][]string
}

func (form *SignUpForm) GetErrors() map[string][]string {
	return form.formErrors
}

func (form *SignUpForm) SetErrors(errs map[string][]string) {
	form.formErrors = errs
}

func CreateSignUpForm(r *http.Request) *SignUpForm {
	r.ParseForm()

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	confirmPassword := r.Form.Get("confirm-password")
	rememberMe := r.Form.Has("remember-me")
	errors := make(map[string][]string)

	return &SignUpForm{
		email:           email,
		password:        password,
		confirmPassword: confirmPassword,
		rememberMe:      rememberMe,
		formErrors:      errors,
	}
}

func (form *SignUpForm) Validate() bool {
	valid := true

	if err := validators.ValidateEmail(form.email); err != nil {
		AppendErrors(err, "email", form)
		valid = false
	}

	if err := validators.ValidatePassword(form.password, 8); err != nil {
		AppendErrors(err, "password", form)
		valid = false
	}

	if err := validators.ValidatePasswordsMatch(form.password, form.confirmPassword); err != nil {
		AppendErrors(err, "confirm_password", form)
		valid = false
	}

	return valid
}

func (form *SignUpForm) ValidateWithoutEmpty() bool {
	valid := true

	if err := validators.ValidateEmailWithoutEmpty(form.email); err != nil {
		AppendErrors(err, "email", form)
		valid = false
	}

	if err := validators.ValidatePasswordWithoutEmpty(form.password, 8); err != nil {
		AppendErrors(err, "password", form)
		valid = false
	}

	if err := validators.ValidatePasswordsMatch(form.password, form.confirmPassword); err != nil {
		AppendErrors(err, "confirm_password", form)
		valid = false
	}

	return valid
}

func (form *SignUpForm) IsValid() bool {
	return len(form.formErrors) == 0
}

func (form *SignUpForm) GetEmail() string {
	return form.email
}

func (form *SignUpForm) GetPassword() string {
	return form.password
}
