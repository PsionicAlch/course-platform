package accounts

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/forms"
	"github.com/PsionicAlch/psionicalch-home/internal/session"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/internal/views"
)

type Handlers struct {
	utils.Loggers
	views   *views.Views
	session session.Session
	auth    *authentication.Authentication
}

func SetupHandlers(views *views.Views, session session.Session, auth *authentication.Authentication) *Handlers {
	loggers := utils.CreateLoggers("ACCOUNT HANDLERS")

	return &Handlers{
		Loggers: loggers,
		views:   views,
		session: session,
		auth:    auth,
	}
}

func (h *Handlers) LoginGet(w http.ResponseWriter, r *http.Request) {
	h.views.RenderHTML(w, "login.page.tmpl", nil)
}

func (h *Handlers) SignupGet(w http.ResponseWriter, r *http.Request) {
	signUpForm := h.session.RetrieveSignUpFormData(r.Context())

	h.views.RenderHTML(w, "signup.page.tmpl", &SignUpPageData{
		SignUpFormData: signUpForm,
	})
}

func (h *Handlers) SignupPost(w http.ResponseWriter, r *http.Request) {
	signUpForm := forms.CreateSignUpForm(r)
	valid := signUpForm.Validate()

	if !valid {
		h.session.StoreSignUpFormData(r.Context(), signUpForm)
		utils.Redirect(w, r, "/accounts/signup")
		return
	}

	err := h.auth.SignUserIn(signUpForm)
	if err != nil {
		switch err.(type) {
		case authentication.UserAlreadyExists:
			forms.AppendError("email", err.Error(), signUpForm)
		}

		h.session.StoreSignUpFormData(r.Context(), signUpForm)
		utils.Redirect(w, r, "/accounts/signup")
		return
	}

	utils.Redirect(w, r, "/")
}

func (h *Handlers) ForgotGet(w http.ResponseWriter, r *http.Request) {
	h.views.RenderHTML(w, "forgot.page.tmpl", nil)
}

func (h *Handlers) ValidateSignupForm(w http.ResponseWriter, r *http.Request) {
	signUpForm := forms.CreateSignUpForm(r)
	signUpForm.ValidateWithoutEmpty()

	h.views.RenderHTMX(w, "signup-form.htmx.tmpl", signUpForm)
}
