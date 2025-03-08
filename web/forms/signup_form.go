package forms

import (
	"net/http"

	"github.com/PsionicAlch/course-platform/web/forms/validators"
	"github.com/PsionicAlch/course-platform/web/html"
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
	firstNameInput.ValidationURL = SignupValidationURL

	lastNameInput := new(html.FormControlComponent)
	lastNameInput.Label = "Last Name:"
	lastNameInput.Name = LastName
	lastNameInput.Type = "text"
	lastNameInput.ValidationURL = SignupValidationURL

	emailInput := new(html.FormControlComponent)
	emailInput.Label = "Email:"
	emailInput.Name = EmailName
	emailInput.Type = "email"
	emailInput.ValidationURL = SignupValidationURL

	passwordInput := new(html.PasswordControlComponent)
	passwordInput.Label = "Password:"
	passwordInput.Name = PasswordName
	passwordInput.ValidationURL = SignupValidationURL

	confirmPasswordInput := new(html.PasswordControlComponent)
	confirmPasswordInput.Label = "Confirm Password:"
	confirmPasswordInput.Name = ConfirmPasswordName
	confirmPasswordInput.ValidationURL = SignupValidationURL

	signupForm := new(html.SignupFormComponent)
	signupForm.FirstNameInput = firstNameInput
	signupForm.LastNameInput = lastNameInput
	signupForm.EmailInput = emailInput
	signupForm.PasswordInput = passwordInput
	signupForm.ConfirmPasswordInput = confirmPasswordInput

	return signupForm
}

func NewSignupFormComponent(form *GenericForm) *html.SignupFormComponent {
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

func GetSignupFormValues(form *GenericForm) (firstName, lastName, email, password, confirmPassword string) {
	firstName = form.GetValue(FirstName)
	lastName = form.GetValue(LastName)
	email = form.GetValue(EmailName)
	password = form.GetValue(PasswordName)
	confirmPassword = form.GetValue(ConfirmPasswordName)

	return
}
