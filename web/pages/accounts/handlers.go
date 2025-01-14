package accounts

import (
	"net"
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/web/forms"
	"github.com/PsionicAlch/psionicalch-home/web/html"
	"github.com/PsionicAlch/psionicalch-home/web/pages"
	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
)

type Handlers struct {
	utils.Loggers
	*pages.HandlerContext
}

func SetupHandlers(handlerContext *pages.HandlerContext) *Handlers {
	loggers := utils.CreateLoggers("ACCOUNT HANDLERS")

	return &Handlers{
		Loggers:        loggers,
		HandlerContext: handlerContext,
	}
}

func (h *Handlers) LoginGet(w http.ResponseWriter, r *http.Request) {
	pageData := html.AccountsLoginPage{
		BasePage:  html.NewBasePage(nil, nosurf.Token(r)),
		LoginForm: forms.EmptyLoginFormComponent(),
	}

	err := h.Renderers.Page.RenderHTML(w, r.Context(), "accounts-login", pageData)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) LoginPost(w http.ResponseWriter, r *http.Request) {
	loginForm := forms.NewLoginForm(r)

	if !loginForm.Validate() {
		loginFormComponent := forms.NewLoginFormComponent(loginForm)
		if err := h.Renderers.Htmx.RenderHTML(w, nil, "login-form", loginFormComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	email, password := forms.GetLoginFormValues(loginForm)
	ipAddr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		h.ErrorLog.Printf("Failed to get IP address from r.RemoteAddr: %s\n", err)

		loginFormComponent := forms.NewLoginFormComponent(loginForm)
		loginFormComponent.ErrorMessage = "Unexpected server error. Please try again."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "login-form", loginFormComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	user, cookie, err := h.Authentication.LogUserIn(email, password)
	if err != nil {
		if err == authentication.ErrInvalidCredentials {
			loginForm.SetEmailError("invalid email or password")
			loginFormComponent := forms.NewLoginFormComponent(loginForm)

			if err := h.Renderers.Htmx.RenderHTML(w, nil, "login-form", loginFormComponent); err != nil {
				h.ErrorLog.Println(err)
			}
		} else {
			h.ErrorLog.Printf("Failed to log user (\"%s\") in: %s\n", email, err)

			loginFormComponent := forms.NewLoginFormComponent(loginForm)
			loginFormComponent.ErrorMessage = "Unexpected server error. Please try again."

			if err := h.Renderers.Htmx.RenderHTML(w, nil, "login-form", loginFormComponent); err != nil {
				h.ErrorLog.Println(err)
			}
		}

		return
	}

	userIpAddresses, err := h.Database.GetUserIpAddresses(user.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to get user's (\"%s\") whitelisted IP addresses: %s\n", user.Email, err)
	} else {
		_, whitelistedIP := utils.InSliceFunc(ipAddr, userIpAddresses, func(ip string, addr *models.WhitelistedIPModel) bool {
			return ip == addr.IPAddress
		})

		if !whitelistedIP {
			go h.Emailer.SendLoginEmail(email, user.Name, ipAddr, time.Now())
		}
	}

	http.SetCookie(w, cookie)

	redirectUrl := h.Session.GetRedirectURL(r.Context())
	if redirectUrl != "" {
		utils.Redirect(w, r, redirectUrl)
	} else {
		utils.Redirect(w, r, "/profile")
	}
}

func (h *Handlers) SignupGet(w http.ResponseWriter, r *http.Request) {
	signupForm := forms.EmptySignupFormComponent()
	pageData := html.AccountsSignupPage{
		BasePage:   html.NewBasePage(nil, nosurf.Token(r)),
		SignupForm: signupForm,
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "accounts-signup", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) SignupPost(w http.ResponseWriter, r *http.Request) {
	signupForm := forms.NewSignupForm(r)

	if !signupForm.Validate() {
		signupFormComponent := forms.NewSignupFormComponent(signupForm)
		if err := h.Renderers.Htmx.RenderHTML(w, nil, "signup-form", signupFormComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	firstName, lastName, email, password, _ := forms.GetSignupFormValues(signupForm)
	ipAddr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		h.ErrorLog.Printf("Failed to get IP address from r.RemoteAddr: %s\n", err)

		signupFormComponent := forms.NewSignupFormComponent(signupForm)
		signupFormComponent.ErrorMessage = "Unexpected server error. Please try again."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "signup-form", signupFormComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	user, cookie, err := h.Authentication.SignUserUp(firstName, lastName, email, password, ipAddr)
	if err != nil {
		if err == authentication.ErrUserExists {
			signupForm.SetEmailError("this email has already been registered")
			signupFormComponent := forms.NewSignupFormComponent(signupForm)

			if err := h.Renderers.Htmx.RenderHTML(w, nil, "signup-form", signupFormComponent); err != nil {
				h.ErrorLog.Println(err)
			}
		} else {
			h.ErrorLog.Printf("Failed to sign user up: %s\n", err)

			signupFormComponent := forms.NewSignupFormComponent(signupForm)
			signupFormComponent.ErrorMessage = "Unexpected server error. Please try again."

			if err := h.Renderers.Htmx.RenderHTML(w, nil, "signup-form", signupFormComponent); err != nil {
				h.ErrorLog.Println(err)
			}
		}

		return
	}

	// TODO: Generate 50% discount code with one usage and add it to their email address.
	go h.Emailer.SendWelcomeEmail(user.Email, user.Name, user.AffiliateCode)

	http.SetCookie(w, cookie)

	redirectUrl := h.Session.GetRedirectURL(r.Context())
	if redirectUrl != "" {
		utils.Redirect(w, r, redirectUrl)
	} else {
		utils.Redirect(w, r, "/courses")
	}
}

func (h *Handlers) LogoutDelete(w http.ResponseWriter, r *http.Request) {
	cookie, err := h.Authentication.LogUserOut(r.Cookies())
	if err != nil {
		h.ErrorLog.Printf("An error occurred whilst logging user out: %s\n", err)
	}

	h.Session.Reset(r.Context())

	http.SetCookie(w, cookie)

	utils.Redirect(w, r, "/")
}

func (h *Handlers) ForgotPasswordGet(w http.ResponseWriter, r *http.Request) {
	forgotPasswordForm := forms.NewForgotPasswordForm(r)
	pageData := html.AccountsForgotPasswordPage{
		BasePage:           html.NewBasePage(nil, nosurf.Token(r)),
		ForgotPasswordForm: forms.NewForgotPasswordFormComponent(forgotPasswordForm),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "accounts-forgot-password", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ForgotPasswordPost(w http.ResponseWriter, r *http.Request) {
	forgotPasswordForm := forms.NewForgotPasswordForm(r)

	if !forgotPasswordForm.Validate() {
		forgotPasswordFormComponent := forms.NewForgotPasswordFormComponent(forgotPasswordForm)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "forgot-password-form", forgotPasswordFormComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	email := forms.GetForgotPasswordFormValues(forgotPasswordForm)
	user, resetToken, err := h.Authentication.GeneratePasswordResetToken(email)
	if err != nil && err != authentication.ErrUnregisteredEmail {
		h.ErrorLog.Printf("Failed to generate new password reset token: %s\n", err)

		forgotPasswordFormComponent := forms.NewForgotPasswordFormComponent(forgotPasswordForm)
		forgotPasswordFormComponent.ErrorMessage = "Unexpected server error. Please try again."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "forgot-password-form", forgotPasswordFormComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err == nil {
		go h.Emailer.SendPasswordResetEmail(email, user.Name, resetToken)
	}

	resp := "<p>An email with instructions on how to reset your password has been sent to your inbox. Please note that if there is no account with this email registered then you won't receive any email.</p>"

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ResetPasswordGet(w http.ResponseWriter, r *http.Request) {
	emailToken := chi.URLParam(r, "email_token")
	pageData := html.AccountsResetPasswordPage{
		BasePage:          html.NewBasePage(nil, nosurf.Token(r)),
		ResetPasswordForm: forms.EmptyResetPasswordFormComponent(emailToken),
	}

	valid, err := h.Authentication.ValidateEmailToken(emailToken)
	if err != nil {
		h.ErrorLog.Printf("Failed to validate password reset token: %s\n", err)

		h.Session.SetErrorMessage(r.Context(), "Unexpected server error. Please try again.")

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "accounts-reset-password", pageData); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if !valid {
		h.Session.SetErrorMessage(r.Context(), "The password token is invalid or expired.")
		utils.Redirect(w, r, "/accounts/reset-password")
		return
	}

	pageData.ResetPasswordForm = forms.EmptyResetPasswordFormComponent(emailToken)
	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "accounts-reset-password", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ResetPasswordPost(w http.ResponseWriter, r *http.Request) {
	emailToken := chi.URLParam(r, "email_token")
	resetPasswordForm := forms.NewResetPasswordForm(r)

	if !resetPasswordForm.Validate() {
		resetPasswordFormComponent := forms.NewResetPasswordFormComponent(resetPasswordForm, emailToken)
		if err := h.Renderers.Htmx.RenderHTML(w, nil, "reset-password-form", resetPasswordFormComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	valid, err := h.Authentication.ValidateEmailToken(emailToken)
	if err != nil {
		h.ErrorLog.Printf("Failed to validate password reset token: %s\n", err)

		resetPasswordFormComponent := forms.NewResetPasswordFormComponent(resetPasswordForm, emailToken)
		resetPasswordFormComponent.ErrorMessage = "Unexpected server error. Please try again."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "reset-password-form", resetPasswordFormComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if !valid {
		h.Session.SetErrorMessage(r.Context(), "The password reset token is invalid or expired.")
		utils.Redirect(w, r, "/accounts/reset-password")
		return
	}

	user, err := h.Authentication.GetUserFromEmailToken(emailToken)
	if err != nil {
		h.ErrorLog.Printf("Failed to get user from password reset token: %s\n", err)

		resetPasswordFormComponent := forms.NewResetPasswordFormComponent(resetPasswordForm, emailToken)
		resetPasswordFormComponent.ErrorMessage = "Unexpected server error. Please try again."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "reset-password-form", resetPasswordFormComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	password, _ := forms.GetResetPasswordFormValues(resetPasswordForm)
	err = h.Authentication.ChangeUserPassword(user, password)
	if err != nil {
		h.ErrorLog.Printf("Failed to get user from password reset token: %s\n", err)

		resetPasswordFormComponent := forms.NewResetPasswordFormComponent(resetPasswordForm, emailToken)
		resetPasswordFormComponent.ErrorMessage = "Unexpected server error. Please try again."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "reset-password-form", resetPasswordFormComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	go h.Emailer.SendPasswordResetConfirmationEmail(user.Email, user.Name)

	h.Session.SetInfoMessage(r.Context(), "Your password has been successfully changed. You can now log in using your new credentials.")

	utils.Redirect(w, r, "/accounts/login")
}

func (h *Handlers) ValidateSignupPost(w http.ResponseWriter, r *http.Request) {
	signupForm := forms.SignupFormPartialValidation(r)
	signupForm.Validate()

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "signup-form", forms.NewSignupFormComponent(signupForm)); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ValidateResetPasswordPost(w http.ResponseWriter, r *http.Request) {
	emailToken := chi.URLParam(r, "email_token")
	resetPasswordForm := forms.ResetPasswordFormPartialValidation(r)
	resetPasswordForm.Validate()

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "reset-password-form", forms.NewResetPasswordFormComponent(resetPasswordForm, emailToken)); err != nil {
		h.ErrorLog.Println(err)
	}
}
