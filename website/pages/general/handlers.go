package general

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
	loggers := utils.CreateLoggers("GENERAL HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		renderers: pages.CreateRenderers(pageRenderer, nil),
		auth:      auth,
		db:        db,
	}
}

func (h *Handlers) HomeGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "home.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) AffiliateProgramGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "affiliate-program.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PrivacyPolicyGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "privacy-policy.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) RefundPolicyGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "refund-policy.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}
