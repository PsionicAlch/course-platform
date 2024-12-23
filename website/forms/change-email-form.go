package forms

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/website/forms/validators"
	"github.com/PsionicAlch/psionicalch-home/website/html"
)

func NewChangeEmailForm(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		EmailName: validators.ChainValidators(
			validators.NotEmpty,
			validators.MaxLength(255),
			validators.IsEmail,
			validators.IsNotDisposableEmail,
		),
	})
}

func ChangeEmailFormPartialValidation(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		EmailName: validators.ChainValidators(
			validators.MaxLength(255),
			validators.IsEmail,
			validators.IsNotDisposableEmail,
		),
	})
}

func EmptyChangeEmailFormComponent() *html.ChangeEmailFormComponent {
	emailInput := new(html.FormControlComponent)
	emailInput.Label = "Email:"
	emailInput.Name = EmailName
	emailInput.Type = "email"
	emailInput.ValidationURL = ""

	changeEmailForm := new(html.ChangeEmailFormComponent)
	changeEmailForm.EmailInput = emailInput

	return changeEmailForm
}

func NewChangeEmailFormComponent(form *GenericForm) *html.ChangeEmailFormComponent {
	changeEmailFormComponent := EmptyChangeEmailFormComponent()

	changeEmailFormComponent.EmailInput.Value = form.GetValue(EmailName)
	changeEmailFormComponent.EmailInput.Errors = form.GetErrors(EmailName)

	return changeEmailFormComponent
}

func GetChangeEmailFormValues(form *GenericForm) (email string) {
	email = form.GetValue(EmailName)

	return
}
