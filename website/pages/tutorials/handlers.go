package tutorials

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
	renderers *pages.Renderers
	db        database.Database
}

func SetupHandlers(pageRenderer render.Renderer, db database.Database) *Handlers {
	loggers := utils.CreateLoggers("TUTORIALS HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		renderers: pages.CreateRenderers(pageRenderer, nil),
		db:        db,
	}
}

// TODO: Set up logic for commenting, liking, and bookmarking.

func (h *Handlers) TutorialsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.TutorialsPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, "tutorials", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.TutorialsTutorialPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, "tutorials-tutorial", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}
