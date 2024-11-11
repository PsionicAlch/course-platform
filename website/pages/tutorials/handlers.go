package tutorials

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/pkg/gatekeeper"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	renderers *pages.Renderers
	auth      *gatekeeper.Gatekeeper
	db        database.Database
}

func SetupHandlers(pageRenderer render.Renderer, auth *gatekeeper.Gatekeeper, db database.Database) *Handlers {
	loggers := utils.CreateLoggers("TUTORIALS HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		renderers: pages.CreateRenderers(pageRenderer, nil),
		auth:      auth,
		db:        db,
	}
}

// TODO: Set up logic for commenting, liking, and bookmarking.

func (h *Handlers) TutorialsGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "tutorials.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println("Failed to render tutorials.page.tmpl: ", err)
	}
}

func (h *Handlers) TutorialGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "tutorial.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println("Failed to render tutorial.page.tmpl: ", err)
	}
}
