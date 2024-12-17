package settings

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
	loggers := utils.CreateLoggers("SETTINGS HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		renderers: *pages.CreateRenderers(pageRenderer, nil, nil),
	}
}

func (h *Handlers) SettingsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.SettingsPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, r.Context(), "settings", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}
