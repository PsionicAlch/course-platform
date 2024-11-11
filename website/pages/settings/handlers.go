package settings

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	renderers pages.Renderers
}

func SetupHandlers(pageRenderer render.Renderer) *Handlers {
	loggers := utils.CreateLoggers("SETTINGS HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		renderers: *pages.CreateRenderers(pageRenderer, nil),
	}
}

func (h *Handlers) SettingsGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "settings.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}
