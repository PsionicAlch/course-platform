package general

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/internal/views"
)

type Handlers struct {
	utils.Loggers
	views *views.Views
}

func SetupHandlers(views *views.Views) *Handlers {
	loggers := utils.CreateLoggers("GENERAL HANDLERS")

	return &Handlers{
		Loggers: loggers,
		views:   views,
	}
}

func (h *Handlers) HomeGet(w http.ResponseWriter, r *http.Request) {
	homePageData := CreateHomePageData()

	h.views.RenderHTML(w, "home.page.tmpl", homePageData)
}
