package general

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/pkg/gatekeeper"
	"github.com/PsionicAlch/psionicalch-home/website/html"
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
	user, err := pages.GetUser(r.Cookies(), h.auth, h.db)
	if err != nil {
		h.WarningLog.Println(err)
	}

	homePageData := html.CreateHomePageData(user)

	h.renderers.Page.RenderHTML(w, "home.page.tmpl", homePageData)
}
