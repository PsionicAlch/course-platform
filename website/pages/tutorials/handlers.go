package tutorials

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
