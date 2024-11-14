package forms

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/website/forms/validators"
	"github.com/PsionicAlch/psionicalch-home/website/html"
)

const (
	FirstName           = "first_name"
	LastName            = "last_name"
	EmailName           = "email"
	PasswordName        = "password"
	ConfirmPasswordName = "confirm_password"
	ValidationURL       = "/accounts/validate/signup"
)

func NewSignupForm(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		FirstName: validators.ChainValidators(
			validators.NotEmpty,
			validators.MaxLength(50),
			validators.Profanities,
			validators.NoSpecialCharacters,
			validators.NoNumbers,
		),
		LastName: validators.ChainValidators(
			validators.NotEmpty,
			validators.MaxLength(50),
			validators.Profanities,
			validators.NoSpecialCharacters,
			validators.NoNumbers,
		),
		EmailName: validators.ChainValidators(
			validators.NotEmpty,
			validators.MaxLength(255),
			validators.IsEmail,
			validators.IsNotDisposableEmail,
		),
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

func SignupFormPartialValidation(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		FirstName: validators.ChainValidators(
			validators.MaxLength(50),
			validators.Profanities,
			validators.NoSpecialCharacters,
			validators.NoNumbers,
		),
		LastName: validators.ChainValidators(
			validators.MaxLength(50),
			validators.Profanities,
			validators.NoSpecialCharacters,
			validators.NoNumbers,
		),
		EmailName: validators.ChainValidators(
			validators.MaxLength(255),
			validators.IsEmail,
			validators.IsNotDisposableEmail,
		),
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

func EmptySignupFormComponent() *html.SignupFormComponent {
	firstNameInput := new(html.FormControlComponent)
	firstNameInput.Label = "First Name:"
	firstNameInput.Name = FirstName
	firstNameInput.Type = "text"
	firstNameInput.ValidationURL = ValidationURL

	lastNameInput := new(html.FormControlComponent)
	lastNameInput.Label = "Last Name:"
	lastNameInput.Name = LastName
	lastNameInput.Type = "text"
	lastNameInput.ValidationURL = ValidationURL

	emailInput := new(html.FormControlComponent)
	emailInput.Label = "Email:"
	emailInput.Name = EmailName
	emailInput.Type = "email"
	emailInput.ValidationURL = ValidationURL

	passwordInput := new(html.PasswordControlComponent)
	passwordInput.Label = "Password:"
	passwordInput.Name = PasswordName
	passwordInput.ValidationURL = ValidationURL

	confirmPasswordInput := new(html.PasswordControlComponent)
	confirmPasswordInput.Label = "Confirm Password:"
	confirmPasswordInput.Name = ConfirmPasswordName
	confirmPasswordInput.ValidationURL = ValidationURL

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

	signupFormComponent.FirstNameInput.Value = form.GetValue(FirstName)
	signupFormComponent.FirstNameInput.Errors = form.GetErrors(FirstName)

	signupFormComponent.LastNameInput.Value = form.GetValue(LastName)
	signupFormComponent.LastNameInput.Errors = form.GetErrors(LastName)

	signupFormComponent.EmailInput.Value = form.GetValue(EmailName)
	signupFormComponent.EmailInput.Errors = form.GetErrors(EmailName)

	signupFormComponent.PasswordInput.Value = form.GetValue(PasswordName)
	signupFormComponent.PasswordInput.Errors = form.GetErrors(PasswordName)

	signupFormComponent.ConfirmPasswordInput.Value = form.GetValue(ConfirmPasswordName)
	signupFormComponent.ConfirmPasswordInput.Errors = form.GetErrors(ConfirmPasswordName)

	return signupFormComponent
}

func GetFormValues(form *GenericForm) (firstName, lastName, email, password, confirmPassword string) {
	firstName = form.GetValue(FirstName)
	lastName = form.GetValue(LastName)
	email = form.GetValue(EmailName)
	password = form.GetValue(PasswordName)
	confirmPassword = form.GetValue(ConfirmPasswordName)

	return
}

func SetEmailError(form *GenericForm, err string) {
	form.SetError(EmailName, err)
}
