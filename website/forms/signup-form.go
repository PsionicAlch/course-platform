package forms

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/website/html"
)

func NewSignupForm(r *http.Request) *html.SignupFormComponent {
	r.ParseForm()

	validationURL := "/accounts/validate/signup"

	firstNameInput := new(html.FormControlComponent)
	firstNameInput.Label = "First Name:"
	firstNameInput.Name = "first-name"
	firstNameInput.Type = "text"
	firstNameInput.ValidationURL = validationURL
	firstNameInput.Value = r.Form.Get(firstNameInput.Name)

	lastNameInput := new(html.FormControlComponent)
	lastNameInput.Label = "Last Name:"
	lastNameInput.Name = "last-name"
	lastNameInput.Type = "text"
	lastNameInput.ValidationURL = validationURL
	lastNameInput.Value = r.Form.Get(lastNameInput.Name)

	emailInput := new(html.FormControlComponent)
	emailInput.Label = "Email:"
	emailInput.Name = "email"
	emailInput.Type = "email"
	emailInput.ValidationURL = validationURL
	emailInput.Value = r.Form.Get(emailInput.Name)

	passwordInput := new(html.PasswordControlComponent)
	passwordInput.Label = "Password:"
	passwordInput.Name = "password"
	passwordInput.ValidationURL = validationURL
	passwordInput.Value = r.Form.Get(passwordInput.Name)

	confirmPasswordInput := new(html.PasswordControlComponent)
	confirmPasswordInput.Label = "Confirm Password:"
	confirmPasswordInput.Name = "confirm_password"
	confirmPasswordInput.ValidationURL = validationURL
	confirmPasswordInput.Value = r.Form.Get(confirmPasswordInput.Name)

	signupForm := new(html.SignupFormComponent)
	signupForm.FirstNameInput = firstNameInput
	signupForm.LastNameInput = lastNameInput
	signupForm.EmailInput = emailInput
	signupForm.PasswordInput = passwordInput
	signupForm.ConfirmPasswordInput = confirmPasswordInput

	return signupForm
}
