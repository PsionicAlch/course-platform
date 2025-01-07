package forms

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/website/forms/validators"
	"github.com/PsionicAlch/psionicalch-home/website/html"
)

func NewForgotPasswordForm(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		EmailName: validators.NotEmpty,
	})
}

func EmptyForgotPasswordFormComponent() *html.ForgotPasswordFormComponent {
	emailInput := new(html.FormControlComponent)
	emailInput.Label = "Email:"
	emailInput.Name = EmailName
	emailInput.Type = "email"

	return &html.ForgotPasswordFormComponent{
		EmailInput: emailInput,
	}
}

func NewForgotPasswordFormComponent(form *GenericForm) *html.ForgotPasswordFormComponent {
	loginFormComponent := EmptyForgotPasswordFormComponent()

	loginFormComponent.EmailInput.Value = form.GetValue(EmailName)
	loginFormComponent.EmailInput.Errors = form.GetErrors(EmailName)

	return loginFormComponent
}

func GetForgotPasswordFormValues(form *GenericForm) (email string) {
	email = form.GetValue(EmailName)

	return
}
