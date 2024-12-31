package settings

import (
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/payments"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/session"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/emails"
	"github.com/PsionicAlch/psionicalch-home/website/forms"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	utils.Loggers
	Render   *pages.Renderers
	Database database.Database
	Session  *session.Session
	Auth     *authentication.Authentication
	Emailer  *emails.Emails
	Payments *payments.Payments
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, db database.Database, sessions *session.Session, auth *authentication.Authentication, emailer *emails.Emails, pay *payments.Payments) *Handlers {
	loggers := utils.CreateLoggers("SETTINGS HANDLERS")

	return &Handlers{
		Loggers:  loggers,
		Render:   pages.CreateRenderers(pageRenderer, htmxRenderer, nil),
		Database: db,
		Session:  sessions,
		Auth:     auth,
		Emailer:  emailer,
		Payments: pay,
	}
}

func (h *Handlers) SettingsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.SettingsPage{
		BasePage:            html.NewBasePage(user),
		ChangeFirstNameForm: forms.EmptyChangeFirstNameFormComponent(),
		ChangeLastNameForm:  forms.EmptyChangeLastNameFormComponent(),
		ChangeEmailForm:     forms.EmptyChangeEmailFormComponent(),
		ChangePasswordForm:  forms.EmptyChangePasswordFormComponent(),
	}

	whitelistedIPAddresses, err := h.Database.GetUserIpAddresses(user.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to get user's (\"%s\") whitelisted IP addresses: %s\n", user.ID, err)

		if err := h.Render.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.IPAddresses = whitelistedIPAddresses

	courses, err := h.Database.GetAllCoursesBoughtByUser(user.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all courses bought by user (\"%s\"): %s\n", user.ID, err)

		if err := h.Render.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	courses = slices.DeleteFunc(courses, func(course *models.CourseModel) bool {
		return time.Now().After(course.UpdatedAt.Add(time.Hour * 24 * 30))
	})

	pageData.Courses = courses

	if err := h.Render.Page.RenderHTML(w, r.Context(), "settings", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) WhitelistIPAddressPost(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	ipAddr := chi.URLParam(r, "ip-address")

	ipAddr, err := url.QueryUnescape(ipAddr)
	if err != nil {
		h.ErrorLog.Printf("Failed to query unescape IP address from URL param: %s\n", err)
		h.Session.SetErrorMessage(r.Context(), "Unexpected server error. Failed to add IP address to whitelist.")

		w.Header().Set("HX-Refresh", "true")
		utils.Redirect(w, r, "/settings#manage-ip-addresses")

		return
	}

	if err := h.Database.AddIPAddress(user.ID, ipAddr); err != nil {
		h.ErrorLog.Printf("Failed to add IP address to user's (\"%s\") whitelist: %s\n", user.ID, err)
		h.Session.SetErrorMessage(r.Context(), "Unexpected server error. Failed to add IP address to whitelist.")
	}

	w.Header().Set("HX-Refresh", "true")
	utils.Redirect(w, r, "/settings#manage-ip-addresses")
}

func (h *Handlers) ChangeFirstNamePost(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	form := forms.NewChangeFirstNameForm(r)

	if !form.Validate() {
		if err := h.Render.Htmx.RenderHTML(w, nil, "change-first-name-form", forms.NewChangeFirstNameFormComponent(form)); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Database.UpdateUserName(user.ID, forms.GetChangeFirstNameFormValues(form), user.Surname); err != nil {
		h.ErrorLog.Printf("Failed to update user's (\"%s\") first name: %s\n", user.ID, err)
		h.Session.SetErrorMessage(r.Context(), "Unexpected server error. Please try again.")

		w.Header().Set("HX-Refresh", "true")
		utils.Redirect(w, r, "/settings#change-first-name")

		return
	}

	h.Session.SetInfoMessage(r.Context(), "Successfully updated your first name.")

	w.Header().Set("HX-Refresh", "true")
	utils.Redirect(w, r, "/settings#change-first-name")
}

func (h *Handlers) ChangeLastNamePost(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	form := forms.NewChangeLastNameForm(r)

	if !form.Validate() {
		if err := h.Render.Htmx.RenderHTML(w, nil, "change-last-name-form", forms.NewChangeLastNameFormComponent(form)); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Database.UpdateUserName(user.ID, user.Name, forms.GetChangeLastNameFormValues(form)); err != nil {
		h.ErrorLog.Printf("Failed to update user's (\"%s\") last name: %s\n", user.ID, err)
		h.Session.SetErrorMessage(r.Context(), "Unexpected server error. Please try again.")

		w.Header().Set("HX-Refresh", "true")
		utils.Redirect(w, r, "/settings#change-last-name")

		return
	}

	h.Session.SetInfoMessage(r.Context(), "Successfully updated your last name.")

	w.Header().Set("HX-Refresh", "true")
	utils.Redirect(w, r, "/settings#change-last-name")
}

func (h *Handlers) ChangeEmailPost(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	form := forms.NewChangeEmailForm(r)

	if !form.Validate() {
		if err := h.Render.Htmx.RenderHTML(w, nil, "change-email-form", forms.NewChangeEmailFormComponent(form)); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Database.UpdateUserEmail(user.ID, forms.GetChangeEmailFormValues(form)); err != nil {
		h.ErrorLog.Printf("Failed to update user's (\"%s\") email: %s\n", user.ID, err)
		h.Session.SetErrorMessage(r.Context(), "Unexpected server error. Please try again.")

		w.Header().Set("HX-Refresh", "true")
		utils.Redirect(w, r, "/settings#change-email")

		return
	}

	h.Session.SetInfoMessage(r.Context(), "Successfully updated your email address.")

	w.Header().Set("HX-Refresh", "true")
	utils.Redirect(w, r, "/settings#change-email")
}

func (h *Handlers) ChangePasswordPost(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	form := forms.NewChangePasswordForm(r)

	if !form.Validate() {
		if err := h.Render.Htmx.RenderHTML(w, nil, "change-password-form", forms.NewChangePasswordFormComponent(form)); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	previousPassword, newPassword := forms.GetChangePasswordFormValues(form)
	valid, err := h.Auth.ValidatePassword(user, previousPassword)
	if err != nil {
		h.ErrorLog.Printf("Failed to validate user's (\"%s\") previous password: %s\n", user.ID, err)
		h.Session.SetErrorMessage(r.Context(), "Unexpected server error. Please try again.")

		w.Header().Set("HX-Refresh", "true")
		utils.Redirect(w, r, "/settings#change-password")

		return
	}

	if !valid {
		forms.SetPreviousPasswordError(form, "invalid password provided")

		if err := h.Render.Htmx.RenderHTML(w, nil, "change-password-form", forms.NewChangePasswordFormComponent(form)); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Auth.ChangeUserPassword(user, newPassword); err != nil {
		h.ErrorLog.Printf("Failed to update user's (\"%s\") password: %s\n", user.ID, err)
		h.Session.SetErrorMessage(r.Context(), "Unexpected server error. Please try again.")

		w.Header().Set("HX-Refresh", "true")
		utils.Redirect(w, r, "/settings#change-password")

		return
	}

	go h.Emailer.SendPasswordResetConfirmationEmail(user.Email, user.Name)

	_, authCookie, err := h.Auth.LogUserIn(user.Email, newPassword)
	if err != nil {
		h.ErrorLog.Printf("Failed to log user (\"%s\") in after updating password: %s\n", user.ID, err)
		h.Session.SetErrorMessage(r.Context(), "Unexpected server error. Please try again.")

		utils.Redirect(w, r, "/accounts/login")

		return
	}

	http.SetCookie(w, authCookie)

	h.Session.SetInfoMessage(r.Context(), "Successfully update your password.")

	w.Header().Set("HX-Refresh", "true")
	utils.Redirect(w, r, "/settings#change-password")
}

func (h *Handlers) IPAddressDelete(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	ipAddrID := chi.URLParam(r, "ip-address-id")

	if err := h.Database.DeleteIPAddress(ipAddrID, user.ID); err != nil {
		h.ErrorLog.Printf("Failed to delete IP address (\"%s\") from the database: %s\n", ipAddrID, err)
		h.Session.SetErrorMessage(r.Context(), "Unexpected server error. Failed to delete IP address.")
	} else {
		h.Session.SetInfoMessage(r.Context(), "Successfully delete IP address.")
	}

	w.Header().Set("HX-Refresh", "true")
	utils.Redirect(w, r, "/settings#change-password")
}

func (h *Handlers) RequestRefundPost(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)

	courseId := chi.URLParam(r, "course-id")
	course, err := h.Database.GetCourseByID(courseId)
	if err != nil {
		h.ErrorLog.Printf("Failed to get course by ID (\"%s\"): %s\n", courseId, err)
		h.Session.SetErrorMessage(r.Context(), "Unexpected server error. Please try again.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.Payments.RequestRefund(user, course); err != nil {
		h.ErrorLog.Printf("Failed to request course (\"%s\") refund for user (\"%s\"): %s\n", course.ID, user.ID, err)
		h.Session.SetErrorMessage(r.Context(), "Unexpected server error. Please try again.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Session.SetInfoMessage(r.Context(), "Refund successfully requested.")
}

func (h *Handlers) AccountDelete(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)

	if err := h.Database.DeleteUser(user.ID); err != nil {
		h.ErrorLog.Printf("Failed to delete user's (\"%s\") account: %s\n", user.ID, err)
		h.Session.SetErrorMessage(r.Context(), "Unexpected server error. Failed to delete your account.")
		w.Header().Set("HX-Refresh", "true")
		utils.Redirect(w, r, "/settings#delete-account")
		return
	}

	go h.Emailer.SendAccountDeletionEmail(user.Email, user.Name)

	utils.Redirect(w, r, "/")
}

func (h *Handlers) ValidateChangePassword(w http.ResponseWriter, r *http.Request) {
	form := forms.ChangePasswordFormPartialValidation(r)
	form.Validate()

	if err := h.Render.Htmx.RenderHTML(w, nil, "change-password-form", forms.NewChangePasswordFormComponent(form)); err != nil {
		h.ErrorLog.Println(err)
	}
}
