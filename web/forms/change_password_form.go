package forms

import (
	"net/http"

	"github.com/PsionicAlch/course-platform/web/forms/validators"
	"github.com/PsionicAlch/course-platform/web/html"
)

func NewChangePasswordForm(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		PreviousPasswordName: validators.NotEmpty,
		NewPasswordName: validators.ChainValidators(
			validators.NotEmpty,
			validators.MinLength(10),
			validators.MaxLength(65),
			validators.UppercaseCharacters,
			validators.LowercaseCharacters,
			validators.NumberCharacters,
		),
	})
}

func ChangePasswordFormPartialValidation(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		PreviousPasswordName: validators.Empty,
		NewPasswordName: validators.ChainValidators(
			validators.MinLength(10),
			validators.MaxLength(65),
			validators.UppercaseCharacters,
			validators.LowercaseCharacters,
			validators.NumberCharacters,
		),
	})
}

func EmptyChangePasswordFormComponent() *html.ChangePasswordFormComponent {
	previousPasswordInput := new(html.PasswordControlComponent)
	previousPasswordInput.Label = "Old Password:"
	previousPasswordInput.Name = PreviousPasswordName
	previousPasswordInput.ValidationURL = ChangePasswordValidationURL

	newPasswordInput := new(html.PasswordControlComponent)
	newPasswordInput.Label = "New Password:"
	newPasswordInput.Name = NewPasswordName
	newPasswordInput.ValidationURL = ChangePasswordValidationURL

	changePasswordForm := new(html.ChangePasswordFormComponent)
	changePasswordForm.PreviousPasswordInput = previousPasswordInput
	changePasswordForm.NewPasswordInput = newPasswordInput

	return changePasswordForm
}

func NewChangePasswordFormComponent(form *GenericForm) *html.ChangePasswordFormComponent {
	changePasswordFormComponent := EmptyChangePasswordFormComponent()

	changePasswordFormComponent.PreviousPasswordInput.Value = form.GetValue(PreviousPasswordName)
	changePasswordFormComponent.PreviousPasswordInput.Errors = form.GetErrors(PreviousPasswordName)

	changePasswordFormComponent.NewPasswordInput.Value = form.GetValue(NewPasswordName)
	changePasswordFormComponent.NewPasswordInput.Errors = form.GetErrors(NewPasswordName)

	return changePasswordFormComponent
}

func GetChangePasswordFormValues(form *GenericForm) (previousPassword, newPassword string) {
	previousPassword = form.GetValue(PreviousPasswordName)
	newPassword = form.GetValue(NewPasswordName)

	return
}

func SetPreviousPasswordError(form *GenericForm, msg string) {
	form.SetError(PreviousPasswordName, msg)
}
