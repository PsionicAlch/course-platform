package tutorials

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/pkg/gatekeeper"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
	"github.com/go-chi/chi/v5"
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

func (h *Handlers) TutorialsGet(w http.ResponseWriter, r *http.Request) {
	user, err := pages.GetUser(r.Cookies(), h.auth, h.db)
	if err != nil {
		h.WarningLog.Println(err)
	}

	tutorials, err := h.db.GetAllTutorials()
	if err != nil {
		h.ErrorLog.Fatalln("Failed to get all tutorials from the database: ", err)
		return
	}

	tutorialsPageData := html.CreateTutorialsPageData(user, tutorials)

	err = h.renderers.Page.RenderHTML(w, "tutorials.page.tmpl", tutorialsPageData)
	if err != nil {
		h.ErrorLog.Println("Failed to render tutorials.page.tmpl: ", err)
	}
}

func (h *Handlers) TutorialGet(w http.ResponseWriter, r *http.Request) {
	user, err := pages.GetUser(r.Cookies(), h.auth, h.db)
	if err != nil {
		h.WarningLog.Println(err)
	}

	slugParam := chi.URLParam(r, "slug")
	tutorialModel, err := h.db.GetTutorialBySlug(slugParam)
	if err != nil {
		h.ErrorLog.Println("Failed to find tutorial by slug: ", err)
		return
	}

	err = h.renderers.Page.RenderHTML(w, "tutorial.page.tmpl", html.TutorialPageData{
		User:     user,
		Tutorial: tutorialModel,
	})
	if err != nil {
		h.ErrorLog.Println("Failed to render tutorial.page.tmpl: ", err)
	}
}
