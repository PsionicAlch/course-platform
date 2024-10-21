package tutorials

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/internal/views"
)

type Handlers struct {
	utils.Loggers
	views     *views.Views
	tutorials *Tutorials
}

func SetupHandlers(views *views.Views, tutorials *Tutorials) *Handlers {
	loggers := utils.CreateLoggers("TUTORIALS HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		views:     views,
		tutorials: tutorials,
	}
}

func (h *Handlers) TutorialsGet(w http.ResponseWriter, r *http.Request) {
	tutorialsPageData := CreateTutorialsPageData(h.tutorials)

	h.views.RenderHTML(w, "tutorials.page.tmpl", tutorialsPageData)
}
