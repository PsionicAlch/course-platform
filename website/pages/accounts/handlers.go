package accounts

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/forms"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	Renderers *pages.Renderers
	Auth      *authentication.Authentication
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, auth *authentication.Authentication) *Handlers {
	loggers := utils.CreateLoggers("ACCOUNT HANDLERS")

	return &Handlers{
		Loggers: loggers,
		Renderers: &pages.Renderers{
			Page: pageRenderer,
			Htmx: htmxRenderer,
		},
		Auth: auth,
	}
}

func (h *Handlers) LoginGet(w http.ResponseWriter, r *http.Request) {
	loginForm := forms.NewLoginForm(r)

	err := h.Renderers.Page.RenderHTML(w, "accounts-login.page.tmpl", html.AccountsLoginPage{
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

	err := h.Renderers.Page.RenderHTML(w, "accounts-signup.page.tmpl", html.AccountsSignupPage{
		SignupForm: signupForm,
	})
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) SignupPost(w http.ResponseWriter, r *http.Request) {
	signupForm := forms.NewSignupForm(r)

	if !signupForm.Validate() {
		err := h.Renderers.Page.RenderHTML(w, "accounts-signup.page.tmpl", html.AccountsSignupPage{
			SignupForm: forms.NewSignupFormComponent(signupForm),
		})
		if err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	firstName, lastName, email, password, _ := forms.GetFormValues(signupForm)
	cookie, err := h.Auth.SignUserUp(firstName, lastName, email, password, r.RemoteAddr)
	if err != nil {
		if err == authentication.ErrUserExists {
			forms.SetEmailError(signupForm, "this email has already been registered")

			err := h.Renderers.Page.RenderHTML(w, "accounts-signup.page.tmpl", html.AccountsSignupPage{
				SignupForm: forms.NewSignupFormComponent(signupForm),
			})
			if err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		h.ErrorLog.Printf("Failed to sign user up: %s\n", err)

		// TODO: Set flash message about unexpected server error.
		err := h.Renderers.Page.RenderHTML(w, "accounts-signup.page.tmpl", html.AccountsSignupPage{
			SignupForm: forms.NewSignupFormComponent(signupForm),
		})
		if err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	http.SetCookie(w, cookie)

	// TODO: Create sessions system so that we can redirect user back to the page that they were on before.

	// Redirect user to courses page so that they can start buying courses.
	http.Redirect(w, r, "/courses", http.StatusFound)
}

func (h *Handlers) ForgotGet(w http.ResponseWriter, r *http.Request) {
	err := h.Renderers.Page.RenderHTML(w, "accounts-forgot-password.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ResetPasswordGet(w http.ResponseWriter, r *http.Request) {
	err := h.Renderers.Page.RenderHTML(w, "accounts-reset-password.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ValidateLoginPost(w http.ResponseWriter, r *http.Request) {
	loginForm := forms.NewLoginForm(r)

	err := h.Renderers.Htmx.RenderHTML(w, "login-form.htmx.tmpl", loginForm)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ValidateSignupPost(w http.ResponseWriter, r *http.Request) {
	signupForm := forms.NewSignupFormComponent(forms.SignupFormPartialValidation(r))

	err := h.Renderers.Htmx.RenderHTML(w, "signup-form.htmx.tmpl", signupForm)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}
