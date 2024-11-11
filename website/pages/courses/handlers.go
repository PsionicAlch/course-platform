package courses

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
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
	err := h.renderers.Page.RenderHTML(w, "courses.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CourseGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "course.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PurchaseCourseGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "course-purchase.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}
