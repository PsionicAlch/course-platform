package tutorials

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	renderers *pages.Renderers
	tutorials *Tutorials
}

func SetupHandlers(pageRenderer render.Renderer, tutorials *Tutorials) *Handlers {
	loggers := utils.CreateLoggers("TUTORIALS HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		renderers: pages.CreateRenderers(pageRenderer, nil),
		tutorials: tutorials,
	}
}

func (h *Handlers) TutorialsGet(w http.ResponseWriter, r *http.Request) {
	tutorialsPageData := CreateTutorialsPageData(h.tutorials)

	h.renderers.Page.RenderHTML(w, "tutorials.page.tmpl", tutorialsPageData)
}
