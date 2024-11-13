package accounts

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/forms"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	renderers *pages.Renderers
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer) *Handlers {
	loggers := utils.CreateLoggers("ACCOUNT HANDLERS")

	return &Handlers{
		Loggers: loggers,
		renderers: &pages.Renderers{
			Page: pageRenderer,
			Htmx: htmxRenderer,
		},
	}
}

func (h *Handlers) LoginGet(w http.ResponseWriter, r *http.Request) {
	loginForm := forms.NewLoginForm(r)

	err := h.renderers.Page.RenderHTML(w, "accounts-login.page.tmpl", html.AccountsLoginPage{
		LoginForm: loginForm,
	})
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) LoginPost(w http.ResponseWriter, r *http.Request) {
	utils.Redirect(w, r, "/")
}

func (h *Handlers) SignupGet(w http.ResponseWriter, r *http.Request) {
	signupForm := forms.EmptySignupFormComponent()

	err := h.renderers.Page.RenderHTML(w, "accounts-signup.page.tmpl", html.AccountsSignupPage{
		SignupForm: signupForm,
	})
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) SignupPost(w http.ResponseWriter, r *http.Request) {
	signupForm := forms.NewSignupForm(r)

	// Validate form. If it's invalid then we send the form back to the user with the errors.
	if !signupForm.Validate() {
		err := h.renderers.Page.RenderHTML(w, "accounts-signup.page.tmpl", html.AccountsSignupPage{
			SignupForm: forms.NewSignupFormComponent(signupForm),
		})
		if err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

}

func (h *Handlers) ForgotGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "accounts-forgot-password.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ResetPasswordGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "accounts-reset-password.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ValidateLoginPost(w http.ResponseWriter, r *http.Request) {
	loginForm := forms.NewLoginForm(r)

	err := h.renderers.Htmx.RenderHTML(w, "login-form.htmx.tmpl", loginForm)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ValidateSignupPost(w http.ResponseWriter, r *http.Request) {
	signupForm := forms.NewSignupFormComponent(forms.SignupFormPartialValidation(r))

	err := h.renderers.Htmx.RenderHTML(w, "signup-form.htmx.tmpl", signupForm)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}
