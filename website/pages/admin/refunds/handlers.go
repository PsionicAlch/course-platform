package refunds

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	Renderers pages.Renderers
	Database  database.Database
	Auth      *authentication.Authentication
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, db database.Database, auth *authentication.Authentication) *Handlers {
	loggers := utils.CreateLoggers("ADMIN HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: *pages.CreateRenderers(pageRenderer, htmxRenderer, nil),
		Database:  db,
		Auth:      auth,
	}
}

func (h *Handlers) RefundsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminRefundsPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-refunds", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}
