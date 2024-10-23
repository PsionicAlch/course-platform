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
	utils.Redirect(w, r, "/tutorials", http.StatusTemporaryRedirect)
	// h.views.RenderHTML(w, "index.page.tmpl", nil)
}
