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
	loginForm := forms.EmptyLoginFormComponent()

	err := h.Renderers.Page.RenderHTML(w, "accounts-login.page.tmpl", html.AccountsLoginPage{
		LoginForm: loginForm,
	})
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) LoginPost(w http.ResponseWriter, r *http.Request) {
	loginForm := forms.NewLoginForm(r)

	if !loginForm.Validate() {
		if err := h.Renderers.Page.RenderHTML(w, "accounts-login.page.tmpl", html.AccountsLoginPage{
			LoginForm: forms.NewLoginFormComponent(loginForm),
		}); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	email, password := forms.GetLoginFormValues(loginForm)
	cookie, err := h.Auth.LogUserIn(email, password, r.RemoteAddr)
	if err != nil {
		if err == authentication.ErrInvalidCredentials {
			loginForm.SetEmailError("invalid email or password")
			if err := h.Renderers.Page.RenderHTML(w, "accounts-login.page.tmpl", html.AccountsLoginPage{
				LoginForm: forms.NewLoginFormComponent(loginForm),
			}); err != nil {
				h.ErrorLog.Println(err)
			}
		} else {
			h.ErrorLog.Printf("Failed to log user (\"%s\") in: %s\n", email, err)

			// TODO: Set up flash message for unexpected server errors.
			if err := h.Renderers.Page.RenderHTML(w, "accounts-login.page.tmpl", html.AccountsLoginPage{
				LoginForm: forms.NewLoginFormComponent(loginForm),
			}); err != nil {
				h.ErrorLog.Println(err)
			}
		}

		return
	}

	http.SetCookie(w, cookie)

	// TODO: Create sessions system so that we can redirect user back to the page that they were on before.

	// In case we weren't redirected to login, redirect user to their profile page.
	http.Redirect(w, r, "/profile", http.StatusFound)
}

func (h *Handlers) SignupGet(w http.ResponseWriter, r *http.Request) {
	signupForm := forms.EmptySignupFormComponent()

	if err := h.Renderers.Page.RenderHTML(w, "accounts-signup.page.tmpl", html.AccountsSignupPage{
		SignupForm: signupForm,
	}); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) SignupPost(w http.ResponseWriter, r *http.Request) {
	signupForm := forms.NewSignupForm(r)

	if !signupForm.Validate() {
		if err := h.Renderers.Page.RenderHTML(w, "accounts-signup.page.tmpl", html.AccountsSignupPage{
			SignupForm: forms.NewSignupFormComponent(signupForm),
		}); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	firstName, lastName, email, password, _ := forms.GetSignupFormValues(signupForm)
	cookie, err := h.Auth.SignUserUp(firstName, lastName, email, password, r.RemoteAddr)
	if err != nil {
		if err == authentication.ErrUserExists {
			signupForm.SetEmailError("this email has already been registered")

			if err := h.Renderers.Page.RenderHTML(w, "accounts-signup.page.tmpl", html.AccountsSignupPage{
				SignupForm: forms.NewSignupFormComponent(signupForm),
			}); err != nil {
				h.ErrorLog.Println(err)
			}
		} else {
			h.ErrorLog.Printf("Failed to sign user up: %s\n", err)

			// TODO: Set flash message about unexpected server error.
			if err := h.Renderers.Page.RenderHTML(w, "accounts-signup.page.tmpl", html.AccountsSignupPage{
				SignupForm: forms.NewSignupFormComponent(signupForm),
			}); err != nil {
				h.ErrorLog.Println(err)
			}
		}

		return
	}

	http.SetCookie(w, cookie)

	// TODO: Create sessions system so that we can redirect user back to the page that they were on before.

	// Redirect user to courses page so that they can start buying courses.
	http.Redirect(w, r, "/courses", http.StatusFound)
}

func (h *Handlers) ForgotGet(w http.ResponseWriter, r *http.Request) {
	if err := h.Renderers.Page.RenderHTML(w, "accounts-forgot-password.page.tmpl", nil); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ResetPasswordGet(w http.ResponseWriter, r *http.Request) {
	if err := h.Renderers.Page.RenderHTML(w, "accounts-reset-password.page.tmpl", nil); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ValidateSignupPost(w http.ResponseWriter, r *http.Request) {
	signupForm := forms.NewSignupFormComponent(forms.SignupFormPartialValidation(r))

	if err := h.Renderers.Htmx.RenderHTML(w, "signup-form.htmx.tmpl", signupForm); err != nil {
		h.ErrorLog.Println(err)
	}
}
