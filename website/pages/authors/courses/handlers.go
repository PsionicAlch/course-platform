package courses

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	Render   pages.Renderers
	Database database.Database
}

func SetupHandlers(pageRenderer, htmxRenderer render.Renderer, db database.Database) *Handlers {
	loggers := utils.CreateLoggers("AUTHOR COURSES HANDLERS")

	return &Handlers{
		Loggers:  loggers,
		Render:   *pages.CreateRenderers(pageRenderer, htmxRenderer),
		Database: db,
	}
}

func (h *Handlers) CoursesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AuthorsCoursesPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Render.Page.RenderHTML(w, r.Context(), "authors-courses", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CoursesPaginationGet(w http.ResponseWriter, r *http.Request) {

}
