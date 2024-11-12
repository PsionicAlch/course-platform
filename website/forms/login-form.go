package forms

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/website/html"
)

func NewLoginForm(r *http.Request) *html.LoginFormComponent {
	r.ParseForm()

	emailInput := new(html.FormControlComponent)
	emailInput.Label = "Email:"
	emailInput.Name = "email"
	emailInput.Type = "email"
	emailInput.ValidationURL = "/accounts/validate/login"
	emailInput.Value = r.Form.Get("email")

	passwordInput := new(html.PasswordControlComponent)
	passwordInput.Name = "password"
	passwordInput.Label = "Password:"
	passwordInput.ValidationURL = "/accounts/validate/login"
	passwordInput.Value = r.Form.Get("password")

	return &html.LoginFormComponent{
		EmailInput:    emailInput,
		PasswordInput: passwordInput,
	}
}
