package forms

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/website/forms/validators"
	"github.com/PsionicAlch/psionicalch-home/website/html"
)

func NewResetPasswordForm(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		PasswordName: validators.ChainValidators(
			validators.NotEmpty,
			validators.MinLength(10),
			validators.MaxLength(65),
			validators.UppercaseCharacters,
			validators.LowercaseCharacters,
			validators.NumberCharacters,
		),
		ConfirmPasswordName: validators.ChainValidators(
			validators.MatchesField(PasswordName, "Password"),
		),
	})
}

func ResetPasswordFormPartialValidation(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		PasswordName: validators.ChainValidators(
			validators.MinLength(10),
			validators.MaxLength(65),
			validators.UppercaseCharacters,
			validators.LowercaseCharacters,
			validators.NumberCharacters,
		),
		ConfirmPasswordName: validators.ChainValidators(
			validators.MatchesField(PasswordName, "Password"),
		),
	})
}

func EmptyResetPasswordFormComponent(emailToken string) *html.ResetPasswordFormComponent {
	passwordInput := new(html.PasswordControlComponent)
	passwordInput.Label = "Password:"
	passwordInput.Name = PasswordName
	passwordInput.ValidationURL = ResetPasswordValidationURL + "/" + emailToken

	confirmPasswordInput := new(html.PasswordControlComponent)
	confirmPasswordInput.Label = "Confirm Password:"
	confirmPasswordInput.Name = ConfirmPasswordName
	confirmPasswordInput.ValidationURL = ResetPasswordValidationURL + "/" + emailToken

	resetPasswordFormComponent := new(html.ResetPasswordFormComponent)
	resetPasswordFormComponent.EmailToken = emailToken
	resetPasswordFormComponent.PasswordInput = passwordInput
	resetPasswordFormComponent.ConfirmPasswordInput = confirmPasswordInput

	return resetPasswordFormComponent
}

func NewResetPasswordFormComponent(form *GenericForm, emailToken string) *html.ResetPasswordFormComponent {
	resetPasswordFormComponent := EmptyResetPasswordFormComponent(emailToken)

	resetPasswordFormComponent.PasswordInput.Value = form.GetValue(PasswordName)
	resetPasswordFormComponent.PasswordInput.Errors = form.GetErrors(PasswordName)

	resetPasswordFormComponent.ConfirmPasswordInput.Value = form.GetValue(ConfirmPasswordName)
	resetPasswordFormComponent.ConfirmPasswordInput.Errors = form.GetErrors(ConfirmPasswordName)

	return resetPasswordFormComponent
}

func GetResetPasswordFormValues(form *GenericForm) (password, confirmPassword string) {
	password = form.GetValue(PasswordName)
	confirmPassword = form.GetValue(ConfirmPasswordName)

	return
}
