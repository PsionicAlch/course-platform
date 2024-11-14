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
	user := authentication.GetUserFromRequest(r)
	pageData := html.AccountsLoginPage{
		BasePage:  html.NewBasePage(user),
		LoginForm: forms.EmptyLoginFormComponent(),
	}

	err := h.Renderers.Page.RenderHTML(w, "accounts-login.page.tmpl", pageData)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) LoginPost(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	loginForm := forms.NewLoginForm(r)
	pageData := html.AccountsLoginPage{
		BasePage: html.NewBasePage(user),
	}

	if !loginForm.Validate() {
		pageData.LoginForm = forms.NewLoginFormComponent(loginForm)
		if err := h.Renderers.Page.RenderHTML(w, "accounts-login.page.tmpl", pageData); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	email, password := forms.GetLoginFormValues(loginForm)
	cookie, err := h.Auth.LogUserIn(email, password, r.RemoteAddr)
	if err != nil {
		if err == authentication.ErrInvalidCredentials {
			loginForm.SetEmailError("invalid email or password")
			pageData.LoginForm = forms.NewLoginFormComponent(loginForm)

			if err := h.Renderers.Page.RenderHTML(w, "accounts-login.page.tmpl", pageData); err != nil {
				h.ErrorLog.Println(err)
			}
		} else {
			h.ErrorLog.Printf("Failed to log user (\"%s\") in: %s\n", email, err)

			pageData.LoginForm = forms.NewLoginFormComponent(loginForm)

			// TODO: Set up flash message for unexpected server errors.
			if err := h.Renderers.Page.RenderHTML(w, "accounts-login.page.tmpl", pageData); err != nil {
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
	user := authentication.GetUserFromRequest(r)
	signupForm := forms.EmptySignupFormComponent()
	pageData := html.AccountsSignupPage{
		BasePage:   html.NewBasePage(user),
		SignupForm: signupForm,
	}

	if err := h.Renderers.Page.RenderHTML(w, "accounts-signup.page.tmpl", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) SignupPost(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	signupForm := forms.NewSignupForm(r)
	pageData := html.AccountsSignupPage{
		BasePage: html.NewBasePage(user),
	}

	if !signupForm.Validate() {
		pageData.SignupForm = forms.NewSignupFormComponent(signupForm)
		if err := h.Renderers.Page.RenderHTML(w, "accounts-signup.page.tmpl", pageData); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	firstName, lastName, email, password, _ := forms.GetSignupFormValues(signupForm)
	cookie, err := h.Auth.SignUserUp(firstName, lastName, email, password, r.RemoteAddr)
	if err != nil {
		if err == authentication.ErrUserExists {
			signupForm.SetEmailError("this email has already been registered")
			pageData.SignupForm = forms.NewSignupFormComponent(signupForm)

			if err := h.Renderers.Page.RenderHTML(w, "accounts-signup.page.tmpl", pageData); err != nil {
				h.ErrorLog.Println(err)
			}
		} else {
			h.ErrorLog.Printf("Failed to sign user up: %s\n", err)

			pageData.SignupForm = forms.NewSignupFormComponent(signupForm)

			// TODO: Set flash message about unexpected server error.
			if err := h.Renderers.Page.RenderHTML(w, "accounts-signup.page.tmpl", pageData); err != nil {
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
	user := authentication.GetUserFromRequest(r)
	pageData := html.AccountsForgotPasswordPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, "accounts-forgot-password.page.tmpl", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ResetPasswordGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AccountsResetPasswordPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, "accounts-reset-password.page.tmpl", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ValidateSignupPost(w http.ResponseWriter, r *http.Request) {
	signupForm := forms.SignupFormPartialValidation(r)
	signupForm.Validate()

	if err := h.Renderers.Htmx.RenderHTML(w, "signup-form.htmx.tmpl", forms.NewSignupFormComponent(signupForm)); err != nil {
		h.ErrorLog.Println(err)
	}
}
