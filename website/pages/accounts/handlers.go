package accounts

import (
	"fmt"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/session"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/pkg/gatekeeper"
	"github.com/PsionicAlch/psionicalch-home/website/forms"
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
	loginForm := h.session.RetrieveLoginFormData(r.Context())

	h.renderers.Page.RenderHTML(w, "login.page.tmpl", &html.LoginPageData{
		LoginForm: &html.LoginFormComponentData{
			Form:  loginForm,
			Error: "",
		},
	})
}

func (h *Handlers) LoginPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	loginForm := forms.CreateLoginForm(r.Form)
	valid := loginForm.Validate()

	fmt.Printf("%#v\n", loginForm)

	if !valid {
		if utils.IsHTMX(r) {
			h.renderers.Htmx.RenderHTML(w, "login-form.htmx.tmpl", &html.LoginFormComponentData{
				Form:  loginForm,
				Error: "",
			})
		} else {
			h.session.StoreLoginFormData(r.Context(), loginForm)
			utils.Redirect(w, r, "/accounts/login")
		}

		return
	}

	cookie, err := h.auth.LogUserIn(loginForm.Email, loginForm.Password, r.RemoteAddr, loginForm.RememberMe)
	if err != nil {
		formErr := ""
		if _, ok := err.(gatekeeper.UserDoesNotExist); ok {
			loginForm.AddError("email", "this email address is unregistered")
		} else {
			h.Loggers.ErrorLog.Println("failed to log user in: ", err)
			formErr = "unexpected server error. please try again"
		}

		if utils.IsHTMX(r) {
			err = h.renderers.Htmx.RenderHTML(w, "login-form.htmx.tmpl", &html.LoginFormComponentData{
				Form:  loginForm,
				Error: formErr,
			})
			if err != nil {
				h.Loggers.WarningLog.Println("Failed to render HTML: ", err)
			}
		} else {
			h.session.StoreLoginFormData(r.Context(), loginForm)
			utils.Redirect(w, r, "/accounts/login")
		}

		return
	}

	http.SetCookie(w, cookie)
	utils.Redirect(w, r, "/profile")
}

func (h *Handlers) SignupGet(w http.ResponseWriter, r *http.Request) {
	signUpForm := h.session.RetrieveSignUpFormData(r.Context())

	h.renderers.Page.RenderHTML(w, "signup.page.tmpl", &html.SignUpPageData{
		SignUpForm: &html.SignupFormComponentData{
			Form:  signUpForm,
			Error: "",
		},
	})
}

func (h *Handlers) SignupPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	signUpForm := forms.CreateSignUpForm(r.Form)
	valid := signUpForm.Validate()

	if !valid {
		if utils.IsHTMX(r) {
			h.renderers.Htmx.RenderHTML(w, "signup-form.htmx.tmpl", &html.SignupFormComponentData{
				Form:  signUpForm,
				Error: "",
			})
		} else {
			h.session.StoreSignUpFormData(r.Context(), signUpForm)
			utils.Redirect(w, r, "/accounts/signup")
		}

		return
	}

	cookie, err := h.auth.SignUserIn(signUpForm.Email, signUpForm.Password, r.RemoteAddr, signUpForm.RememberMe)
	if err != nil {
		formErr := ""
		if _, ok := err.(gatekeeper.UserAlreadyExists); ok {
			signUpForm.AddError("email", "this email address is already in use")
		} else {
			h.Loggers.ErrorLog.Println("failed to sign user in: ", err)
			formErr = "unexpected server error. please try again"
		}

		if utils.IsHTMX(r) {
			err = h.renderers.Htmx.RenderHTML(w, "signup-form.htmx.tmpl", &html.SignupFormComponentData{
				Form:  signUpForm,
				Error: formErr,
			})
			if err != nil {
				h.Loggers.WarningLog.Println("Failed to render HTML: ", err)
			}
		} else {
			h.session.StoreSignUpFormData(r.Context(), signUpForm)
			utils.Redirect(w, r, "/accounts/signup")
		}

		return
	}

	http.SetCookie(w, cookie)
	utils.Redirect(w, r, "/profile")
}

func (h *Handlers) ForgotGet(w http.ResponseWriter, r *http.Request) {
	h.renderers.Page.RenderHTML(w, "forgot.page.tmpl", nil)
}

func (h *Handlers) ValidateSignupForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	signUpForm := forms.CreateSignUpForm(r.Form)
	signUpForm.ValidateWithoutEmpty()

	h.renderers.Htmx.RenderHTML(w, "signup-form.htmx.tmpl", &html.SignupFormComponentData{
		Form:  signUpForm,
		Error: "",
	})
}
