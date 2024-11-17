package accounts

import (
	"net"
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/emails"
	"github.com/PsionicAlch/psionicalch-home/website/forms"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	utils.Loggers
	Renderers *pages.Renderers
	Auth      *authentication.Authentication
	Emailer   *emails.Emails
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, auth *authentication.Authentication, emailer *emails.Emails) *Handlers {
	loggers := utils.CreateLoggers("ACCOUNT HANDLERS")

	return &Handlers{
		Loggers: loggers,
		Renderers: &pages.Renderers{
			Page: pageRenderer,
			Htmx: htmxRenderer,
		},
		Auth:    auth,
		Emailer: emailer,
	}
}

func (h *Handlers) LoginGet(w http.ResponseWriter, r *http.Request) {
	pageData := html.AccountsLoginPage{
		BasePage:  html.NewBasePage(nil),
		LoginForm: forms.EmptyLoginFormComponent(),
	}

	err := h.Renderers.Page.RenderHTML(w, "accounts-login", pageData)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) LoginPost(w http.ResponseWriter, r *http.Request) {
	loginForm := forms.NewLoginForm(r)
	pageData := html.AccountsLoginPage{
		BasePage: html.NewBasePage(nil),
	}

	if !loginForm.Validate() {
		pageData.LoginForm = forms.NewLoginFormComponent(loginForm)
		if err := h.Renderers.Page.RenderHTML(w, "accounts-login", pageData); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	email, password := forms.GetLoginFormValues(loginForm)
	ipAddr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		h.ErrorLog.Printf("Failed to get IP address from r.RemoteAddr: %s\n", err)

		// TODO: Set up flash message for unexpected server errors.

		pageData.LoginForm = forms.NewLoginFormComponent(loginForm)
		if err := h.Renderers.Page.RenderHTML(w, "accounts-login", pageData); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	user, cookie, err := h.Auth.LogUserIn(email, password)
	if err != nil {
		if err == authentication.ErrInvalidCredentials {
			loginForm.SetEmailError("invalid email or password")
			pageData.LoginForm = forms.NewLoginFormComponent(loginForm)

			if err := h.Renderers.Page.RenderHTML(w, "accounts-login", pageData); err != nil {
				h.ErrorLog.Println(err)
			}
		} else {
			h.ErrorLog.Printf("Failed to log user (\"%s\") in: %s\n", email, err)

			// TODO: Set up flash message for unexpected server errors.

			pageData.LoginForm = forms.NewLoginFormComponent(loginForm)
			if err := h.Renderers.Page.RenderHTML(w, "accounts-login", pageData); err != nil {
				h.ErrorLog.Println(err)
			}
		}

		return
	}

	userIpAddresses, err := h.Auth.Database.GetUserIpAddresses(user.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to get user's (\"%s\") whitelisted IP addresses: %s\n", user.Email, err)
	}

	if userIpAddresses != nil && !utils.InSlice(ipAddr, userIpAddresses) {
		go h.Emailer.SendLoginEmail(email, user.Name, ipAddr, time.Now())
	}

	http.SetCookie(w, cookie)

	// TODO: Create sessions system so that we can redirect user back to the page that they were on before.

	// In case we weren't redirected to login, redirect user to their profile page.
	http.Redirect(w, r, "/profile", http.StatusFound)
}

func (h *Handlers) SignupGet(w http.ResponseWriter, r *http.Request) {
	signupForm := forms.EmptySignupFormComponent()
	pageData := html.AccountsSignupPage{
		BasePage:   html.NewBasePage(nil),
		SignupForm: signupForm,
	}

	if err := h.Renderers.Page.RenderHTML(w, "accounts-signup", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) SignupPost(w http.ResponseWriter, r *http.Request) {
	signupForm := forms.NewSignupForm(r)
	pageData := html.AccountsSignupPage{
		BasePage: html.NewBasePage(nil),
	}

	if !signupForm.Validate() {
		pageData.SignupForm = forms.NewSignupFormComponent(signupForm)
		if err := h.Renderers.Page.RenderHTML(w, "accounts-signup", pageData); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	firstName, lastName, email, password, _ := forms.GetSignupFormValues(signupForm)
	ipAddr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		h.ErrorLog.Printf("Failed to get IP address from r.RemoteAddr: %s\n", err)

		// TODO: Set flash message about unexpected server error.

		pageData.SignupForm = forms.NewSignupFormComponent(signupForm)
		if err := h.Renderers.Page.RenderHTML(w, "accounts-signup", pageData); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	user, cookie, err := h.Auth.SignUserUp(firstName, lastName, email, password, ipAddr)
	if err != nil {
		if err == authentication.ErrUserExists {
			signupForm.SetEmailError("this email has already been registered")
			pageData.SignupForm = forms.NewSignupFormComponent(signupForm)

			if err := h.Renderers.Page.RenderHTML(w, "accounts-signup", pageData); err != nil {
				h.ErrorLog.Println(err)
			}
		} else {
			h.ErrorLog.Printf("Failed to sign user up: %s\n", err)

			// TODO: Set flash message about unexpected server error.

			pageData.SignupForm = forms.NewSignupFormComponent(signupForm)
			if err := h.Renderers.Page.RenderHTML(w, "accounts-signup", pageData); err != nil {
				h.ErrorLog.Println(err)
			}
		}

		return
	}

	go h.Emailer.SendWelcomeEmail(user.Email, user.Name, user.AffiliateCode)

	http.SetCookie(w, cookie)

	// TODO: Create sessions system so that we can redirect user back to the page that they were on before.

	// Redirect user to courses page so that they can start buying courses.
	utils.Redirect(w, r, "/courses")
}

func (h *Handlers) LogoutDelete(w http.ResponseWriter, r *http.Request) {
	h.InfoLog.Printf("Logging user (%#v) out", authentication.GetUserFromRequest(r))

	cookie, err := h.Auth.LogUserOut(r.Cookies())
	if err != nil {
		h.ErrorLog.Printf("An error occurred whilst logging user out: %s\n", err)
	}

	// TODO: Reset session.

	http.SetCookie(w, cookie)

	utils.Redirect(w, r, "/")
}

func (h *Handlers) ForgotPasswordGet(w http.ResponseWriter, r *http.Request) {
	forgotPasswordForm := forms.NewForgotPasswordForm(r)
	pageData := html.AccountsForgotPasswordPage{
		BasePage:           html.NewBasePage(nil),
		ForgotPasswordForm: forms.NewForgotPasswordFormComponent(forgotPasswordForm),
	}

	if err := h.Renderers.Page.RenderHTML(w, "accounts-forgot-password", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ForgotPasswordPost(w http.ResponseWriter, r *http.Request) {
	forgotPasswordForm := forms.NewForgotPasswordForm(r)
	pageData := html.AccountsForgotPasswordPage{
		BasePage:           html.NewBasePage(nil),
		ForgotPasswordForm: nil,
	}

	email := forms.GetForgotPasswordFormValues(forgotPasswordForm)

	user, resetToken, err := h.Auth.GeneratePasswordResetToken(email)
	if err != nil && err != authentication.ErrUnregisteredEmail {
		h.ErrorLog.Printf("Failed to generate new password reset token: %s\n", err)

		// TODO: Set flash message about unexpected server error.

		pageData.ForgotPasswordForm = forms.NewForgotPasswordFormComponent(forgotPasswordForm)
		if err := h.Renderers.Page.RenderHTML(w, "accounts-forgot-password", pageData); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err == nil {
		go h.Emailer.SendPasswordResetEmail(email, user.Name, resetToken)
	}

	if err := h.Renderers.Page.RenderHTML(w, "accounts-forgot-password", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ResetPasswordGet(w http.ResponseWriter, r *http.Request) {
	emailToken := chi.URLParam(r, "email_token")
	pageData := html.AccountsResetPasswordPage{
		BasePage:          html.NewBasePage(nil),
		ResetPasswordForm: forms.EmptyResetPasswordFormComponent(emailToken),
	}

	valid, err := h.Auth.ValidateEmailToken(emailToken)
	if err != nil {
		h.ErrorLog.Printf("Failed to validate password reset token: %s\n", err)

		// TODO: Set flash message about unexpected server error.

		if err := h.Renderers.Page.RenderHTML(w, "accounts-reset-password", pageData); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if !valid {
		// TODO: Set flash message about invalid or expired token.
		utils.Redirect(w, r, "/accounts/reset-password")
		return
	}

	pageData.ResetPasswordForm = forms.EmptyResetPasswordFormComponent(emailToken)
	if err := h.Renderers.Page.RenderHTML(w, "accounts-reset-password", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ResetPasswordPost(w http.ResponseWriter, r *http.Request) {
	emailToken := chi.URLParam(r, "email_token")
	resetPasswordForm := forms.NewResetPasswordForm(r)
	pageData := html.AccountsResetPasswordPage{
		BasePage: html.NewBasePage(nil),
	}

	if !resetPasswordForm.Validate() {
		pageData.ResetPasswordForm = forms.NewResetPasswordFormComponent(resetPasswordForm, emailToken)
		if err := h.Renderers.Page.RenderHTML(w, "accounts-reset-password", pageData); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	valid, err := h.Auth.ValidateEmailToken(emailToken)
	if err != nil {
		h.ErrorLog.Printf("Failed to validate password reset token: %s\n", err)

		// TODO: Set flash message about unexpected server error.

		pageData.ResetPasswordForm = forms.NewResetPasswordFormComponent(resetPasswordForm, emailToken)
		if err := h.Renderers.Page.RenderHTML(w, "accounts-reset-password", pageData); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if !valid {
		// TODO: Set flash message about invalid or expired token.
		utils.Redirect(w, r, "/accounts/reset-password")
		return
	}

	user, err := h.Auth.GetUserFromEmailToken(emailToken)
	if err != nil {
		h.ErrorLog.Printf("Failed to get user from password reset token: %s\n", err)

		// TODO: Set flash message about unexpected server error.

		pageData.ResetPasswordForm = forms.NewResetPasswordFormComponent(resetPasswordForm, emailToken)
		if err := h.Renderers.Page.RenderHTML(w, "accounts-reset-password", pageData); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	password, _ := forms.GetResetPasswordFormValues(resetPasswordForm)
	err = h.Auth.ChangeUserPassword(user, password)
	if err != nil {
		h.ErrorLog.Printf("Failed to get user from password reset token: %s\n", err)

		// TODO: Set flash message about unexpected server error.

		pageData.ResetPasswordForm = forms.NewResetPasswordFormComponent(resetPasswordForm, emailToken)
		if err := h.Renderers.Page.RenderHTML(w, "accounts-reset-password", pageData); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	go h.Emailer.SendPasswordResetConfirmationEmail(user.Email, user.Name)

	// TODO: Set up flash message to inform the user that their password has been changed.

	utils.Redirect(w, r, "/accounts/login")
}

func (h *Handlers) ValidateSignupPost(w http.ResponseWriter, r *http.Request) {
	signupForm := forms.SignupFormPartialValidation(r)
	signupForm.Validate()

	if err := h.Renderers.Htmx.RenderHTML(w, "signup-form", forms.NewSignupFormComponent(signupForm)); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ValidateResetPasswordPost(w http.ResponseWriter, r *http.Request) {
	emailToken := chi.URLParam(r, "email_token")
	resetPasswordForm := forms.ResetPasswordFormPartialValidation(r)
	resetPasswordForm.Validate()

	if err := h.Renderers.Htmx.RenderHTML(w, "reset-password-form", forms.NewResetPasswordFormComponent(resetPasswordForm, emailToken)); err != nil {
		h.ErrorLog.Println(err)
	}
}
