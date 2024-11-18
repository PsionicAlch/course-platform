package courses

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
	renderers *pages.Renderers
}

func SetupHandlers(pageRenderer render.Renderer) *Handlers {
	loggers := utils.CreateLoggers("COURSE HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		renderers: pages.CreateRenderers(pageRenderer, nil),
	}
}

func (h *Handlers) CoursesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.CoursesPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, r.Context(), "courses", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CourseGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.CoursesCoursePage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, r.Context(), "courses-course", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PurchaseCourseGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.CoursesPurchasesPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, r.Context(), "courses-purchase", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}
