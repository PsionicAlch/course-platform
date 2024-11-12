package accounts

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/session"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/pkg/gatekeeper"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	renderers *pages.Renderers
	session   session.Session
	auth      *gatekeeper.Gatekeeper
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, session session.Session, auth *gatekeeper.Gatekeeper) *Handlers {
	loggers := utils.CreateLoggers("ACCOUNT HANDLERS")

	return &Handlers{
		Loggers: loggers,
		renderers: &pages.Renderers{
			Page: pageRenderer,
			Htmx: htmxRenderer,
		},
		session: session,
		auth:    auth,
	}
}

func (h *Handlers) LoginGet(w http.ResponseWriter, r *http.Request) {
	loginForm := new(html.LoginFormComponent)
	loginForm.Email = "email@me.com"
	loginForm.EmailErrors = append(loginForm.EmailErrors, "This is error 1", "This is error 2")

	err := h.renderers.Page.RenderHTML(w, "accounts-login.page.tmpl", html.AccountsLoginPage{
		LoginForm: loginForm,
	})
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) LoginPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Redirect", "/")
}

func (h *Handlers) SignupGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "accounts-signup.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
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
	loginForm := new(html.LoginFormComponent)
	loginForm.Email = "thisemailisvalid@emails.com"
	loginForm.Password = "helloworld123"
	loginForm.PasswordErrors = append(loginForm.PasswordErrors, "This is an error")

	err := h.renderers.Htmx.RenderHTML(w, "login-form.htmx.tmpl", loginForm)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}
