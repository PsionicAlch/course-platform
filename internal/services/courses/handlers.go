package courses

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
	loggers := utils.CreateLoggers("COURSE HANDLERS")

	return &Handlers{
		Loggers: loggers,
		views:   views,
	}
}

func (h *Handlers) CoursesGet(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/tutorials", http.StatusTemporaryRedirect)
	// h.views.RenderHTML(w, "index.page.tmpl", nil)
}
