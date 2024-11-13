package forms

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/website/forms/validators"
	"github.com/PsionicAlch/psionicalch-home/website/html"
)

const (
	firstName           = "first_name"
	lastName            = "last_name"
	emailName           = "email"
	passwordName        = "password"
	confirmPasswordName = "confirm_password"
	validationURL       = "/accounts/validate/signup"
)

func NewSignupForm(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		firstName: validators.ChainValidators(
			validators.NotEmpty,
			validators.MaxLength(50),
			validators.Profanities,
			validators.NoSpecialCharacters,
			validators.NoNumbers,
		),
		lastName: validators.ChainValidators(
			validators.NotEmpty,
			validators.MaxLength(50),
			validators.Profanities,
			validators.NoSpecialCharacters,
			validators.NoNumbers,
		),
		emailName: validators.ChainValidators(
			validators.NotEmpty,
			validators.MaxLength(255),
			validators.IsEmail,
			validators.IsNotDisposableEmail,
		),
		passwordName: validators.ChainValidators(
			validators.NotEmpty,
			validators.MinLength(10),
			validators.MaxLength(65),
			validators.UppercaseCharacters,
			validators.LowercaseCharacters,
			validators.NumberCharacters,
		),
		confirmPasswordName: validators.ChainValidators(
			validators.MatchesField(passwordName, "Password"),
		),
	})
}

func SignupFormPartialValidation(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		firstName: validators.ChainValidators(
			validators.MaxLength(50),
			validators.Profanities,
			validators.NoSpecialCharacters,
			validators.NoNumbers,
		),
		lastName: validators.ChainValidators(
			validators.MaxLength(50),
			validators.Profanities,
			validators.NoSpecialCharacters,
			validators.NoNumbers,
		),
		emailName: validators.ChainValidators(
			validators.MaxLength(255),
			validators.IsEmail,
			validators.IsNotDisposableEmail,
		),
		passwordName: validators.ChainValidators(
			validators.MinLength(10),
			validators.MaxLength(65),
			validators.UppercaseCharacters,
			validators.LowercaseCharacters,
			validators.NumberCharacters,
		),
		confirmPasswordName: validators.ChainValidators(
			validators.MatchesField(passwordName, "Password"),
		),
	})
}

func EmptySignupFormComponent() *html.SignupFormComponent {
	firstNameInput := new(html.FormControlComponent)
	firstNameInput.Label = "First Name:"
	firstNameInput.Name = firstName
	firstNameInput.Type = "text"
	firstNameInput.ValidationURL = validationURL

	lastNameInput := new(html.FormControlComponent)
	lastNameInput.Label = "Last Name:"
	lastNameInput.Name = lastName
	lastNameInput.Type = "text"
	lastNameInput.ValidationURL = validationURL

	emailInput := new(html.FormControlComponent)
	emailInput.Label = "Email:"
	emailInput.Name = emailName
	emailInput.Type = "email"
	emailInput.ValidationURL = validationURL

	passwordInput := new(html.PasswordControlComponent)
	passwordInput.Label = "Password:"
	passwordInput.Name = passwordName
	passwordInput.ValidationURL = validationURL

	confirmPasswordInput := new(html.PasswordControlComponent)
	confirmPasswordInput.Label = "Confirm Password:"
	confirmPasswordInput.Name = confirmPasswordName
	confirmPasswordInput.ValidationURL = validationURL

	signupForm := new(html.SignupFormComponent)
	signupForm.FirstNameInput = firstNameInput
	signupForm.LastNameInput = lastNameInput
	signupForm.EmailInput = emailInput
	signupForm.PasswordInput = passwordInput
	signupForm.ConfirmPasswordInput = confirmPasswordInput

	return signupForm
}

func NewSignupFormComponent(form *GenericForm) *html.SignupFormComponent {
	form.Validate()

	signupFormComponent := EmptySignupFormComponent()

	signupFormComponent.FirstNameInput.Value = form.GetValue("first_name")
	signupFormComponent.FirstNameInput.Errors = form.GetErrors("first_name")

	signupFormComponent.LastNameInput.Value = form.GetValue("last_name")
	signupFormComponent.LastNameInput.Errors = form.GetErrors("last_name")

	signupFormComponent.EmailInput.Value = form.GetValue("email")
	signupFormComponent.EmailInput.Errors = form.GetErrors("email")

	signupFormComponent.PasswordInput.Value = form.GetValue("password")
	signupFormComponent.PasswordInput.Errors = form.GetErrors("password")

	signupFormComponent.ConfirmPasswordInput.Value = form.GetValue("confirm_password")
	signupFormComponent.ConfirmPasswordInput.Errors = form.GetErrors("confirm_password")

	return signupFormComponent
}
