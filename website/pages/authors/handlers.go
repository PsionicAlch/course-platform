package authors

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	renderers pages.Renderers
}

func SetupHandlers(pageRenderer render.Renderer) *Handlers {
	loggers := utils.CreateLoggers("AUTHOR HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		renderers: *pages.CreateRenderers(pageRenderer, nil),
	}
}

func (h *Handlers) TutorialsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AuthorsTutorialsPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, r.Context(), "authors-tutorials", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CoursesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AuthorsCoursesPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, r.Context(), "authors-courses", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}
