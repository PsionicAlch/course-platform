package forms

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/website/forms/validators"
	"github.com/PsionicAlch/psionicalch-home/website/html"
)

func NewLoginForm(r *http.Request) *GenericForm {
	// We're not filling it in because the only real validation
	// will be whether or not the credentials match.
	return NewForm(r, map[FieldName]validators.ValidationFunc{})
}

func EmptyLoginFormComponent() *html.LoginFormComponent {
	emailInput := new(html.FormControlComponent)
	emailInput.Label = "Email:"
	emailInput.Name = EmailName
	emailInput.Type = "email"

	passwordInput := new(html.PasswordControlComponent)
	passwordInput.Name = PasswordName
	passwordInput.Label = "Password:"

	return &html.LoginFormComponent{
		EmailInput:    emailInput,
		PasswordInput: passwordInput,
	}
}

func NewLoginFormComponent(form *GenericForm) *html.LoginFormComponent {
	loginFormComponent := EmptyLoginFormComponent()

	loginFormComponent.EmailInput.Value = form.GetValue(EmailName)
	loginFormComponent.EmailInput.Errors = form.GetErrors(EmailName)

	loginFormComponent.PasswordInput.Value = form.GetValue(PasswordName)
	loginFormComponent.PasswordInput.Errors = form.GetErrors(PasswordName)

	return loginFormComponent
}

func GetLoginFormValues(form *GenericForm) (email, password string) {
	email = form.GetValue(EmailName)
	password = form.GetValue(PasswordName)

	return
}
