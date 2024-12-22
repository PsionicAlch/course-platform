package settings

import (
	"net/http"
	"slices"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/forms"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	Render   *pages.Renderers
	Database database.Database
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, db database.Database) *Handlers {
	loggers := utils.CreateLoggers("SETTINGS HANDLERS")

	return &Handlers{
		Loggers:  loggers,
		Render:   pages.CreateRenderers(pageRenderer, htmxRenderer, nil),
		Database: db,
	}
}

func (h *Handlers) SettingsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.SettingsPage{
		BasePage:            html.NewBasePage(user),
		ChangeFirstNameForm: forms.EmptyChangeFirstNameFormComponent(),
		ChangeLastNameForm:  forms.EmptyChangeLastNameFormComponent(),
		ChangeEmailForm:     forms.EmptyChangeEmailFormComponent(),
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

func (h *Handlers) ValidateChangeFirstName(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) ValidateChangeLastName(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) ValidateChangeEmail(w http.ResponseWriter, r *http.Request) {

}
