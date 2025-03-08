package admin

import (
	"net/http"

	"github.com/PsionicAlch/course-platform/internal/utils"
	"github.com/PsionicAlch/course-platform/web/pages"
)

// TODO: Reset Cache on POST requests

type Handlers struct {
	utils.Loggers
	*pages.HandlerContext
}

func SetupHandlers(handlerContext *pages.HandlerContext) *Handlers {
	loggers := utils.CreateLoggers("ADMIN HANDLERS")

	return &Handlers{
		Loggers:        loggers,
		HandlerContext: handlerContext,
	}
}

func (h *Handlers) AdminGet(w http.ResponseWriter, r *http.Request) {
	utils.Redirect(w, r, "/admin/users")
}
