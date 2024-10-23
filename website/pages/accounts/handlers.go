package accounts

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/forms"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/session"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	renderers *pages.Renderers
	session   session.Session
	auth      *authentication.Authentication
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, session session.Session, auth *authentication.Authentication) *Handlers {
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
	h.renderers.Page.RenderHTML(w, "login.page.tmpl", nil)
}

func (h *Handlers) SignupGet(w http.ResponseWriter, r *http.Request) {
	signUpForm := h.session.RetrieveSignUpFormData(r.Context())

	h.renderers.Page.RenderHTML(w, "signup.page.tmpl", &SignUpPageData{
		SignUpFormData: signUpForm,
	})
}

func (h *Handlers) SignupPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	signUpForm := forms.CreateSignUpForm(r.Form)
	valid := signUpForm.Validate()

	if !valid {
		h.session.StoreSignUpFormData(r.Context(), signUpForm)
		utils.Redirect(w, r, "/accounts/signup")
		return
	}

	cookie, err := h.auth.SignUserIn(signUpForm, r.RemoteAddr)
	if err != nil {
		h.Loggers.ErrorLog.Println("failed to sign user in: ", err)

		h.session.StoreSignUpFormData(r.Context(), signUpForm)
		utils.Redirect(w, r, "/accounts/signup")
		return
	}

	http.SetCookie(w, cookie)
	utils.Redirect(w, r, "/")
}

func (h *Handlers) ForgotGet(w http.ResponseWriter, r *http.Request) {
	h.renderers.Page.RenderHTML(w, "forgot.page.tmpl", nil)
}

func (h *Handlers) ValidateSignupForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	signUpForm := forms.CreateSignUpForm(r.Form)
	signUpForm.ValidateWithoutEmpty()

	h.renderers.Htmx.RenderHTML(w, "signup-form.htmx.tmpl", signUpForm)
}
