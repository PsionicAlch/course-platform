package authors

import (
	"fmt"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Render   pages.Renderers
	Database database.Database
}

func SetupHandlers(pageRenderer, htmxRenderer render.Renderer, db database.Database) *Handlers {
	return &Handlers{
		Render:   *pages.CreateRenderers(pageRenderer, htmxRenderer),
		Database: db,
	}
}

func (h *Handlers) AuthorGet(w http.ResponseWriter, r *http.Request) {
	authorSlug := chi.URLParam(r, "author-slug")
	utils.Redirect(w, r, fmt.Sprintf("/authors/%s/tutorials", authorSlug))
}
