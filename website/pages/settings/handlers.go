package settings

import (
	"net/http"
	"slices"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/session"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/forms"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	Render   *pages.Renderers
	Database database.Database
	Session  *session.Session
	Auth     *authentication.Authentication
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, db database.Database, sessions *session.Session, auth *authentication.Authentication) *Handlers {
	loggers := utils.CreateLoggers("SETTINGS HANDLERS")

	return &Handlers{
		Loggers:  loggers,
		Render:   pages.CreateRenderers(pageRenderer, htmxRenderer, nil),
		Database: db,
		Session:  sessions,
		Auth:     auth,
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

	// TODO: Validate that the old password is valid.
	previousPassword, newPassword := forms.GetChangePasswordFormValues(form)

	if err := h.Auth.ChangeUserPassword(user, newPassword); err != nil {
		// TODO: Error handling.
	}

	_, authCookie, err := h.Auth.LogUserIn(user.Email, newPassword)
	if err != nil {
		// TODO: Error handling.
	}

	http.SetCookie(w, authCookie)

	h.Session.SetInfoMessage(r.Context(), "Successfully update your password.")

	w.Header().Set("HX-Refresh", "true")
	utils.Redirect(w, r, "/settings#change-password")
}

func (h *Handlers) ValidateChangePassword(w http.ResponseWriter, r *http.Request) {
	form := forms.ChangePasswordFormPartialValidation(r)
	form.Validate()

	if err := h.Render.Htmx.RenderHTML(w, nil, "change-password-form", forms.NewChangePasswordFormComponent(form)); err != nil {
		h.ErrorLog.Println(err)
	}
}
