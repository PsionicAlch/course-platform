package profile

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/pkg/gatekeeper"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	renderers pages.Renderers
	auth      *gatekeeper.Gatekeeper
}

func SetupHandlers(pageRenderer render.Renderer, auth *gatekeeper.Gatekeeper) *Handlers {
	loggers := utils.CreateLoggers("PROFILE HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		renderers: *pages.CreateRenderers(pageRenderer, nil),
		auth:      auth,
	}
}

func (h *Handlers) ProfileGet(w http.ResponseWriter, r *http.Request) {
	h.renderers.Page.RenderHTML(w, "profile.page.tmpl", html.ProfilePageData{
		Email: "hello@me.com",
	})
}
