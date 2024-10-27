package profile

import (
	"fmt"
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
	renderers pages.Renderers
	auth      *gatekeeper.Gatekeeper
	db        database.Database
}

func SetupHandlers(pageRenderer render.Renderer, auth *gatekeeper.Gatekeeper, db database.Database) *Handlers {
	loggers := utils.CreateLoggers("PROFILE HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		renderers: *pages.CreateRenderers(pageRenderer, nil),
		auth:      auth,
		db:        db,
	}
}

func (h *Handlers) ProfileGet(w http.ResponseWriter, r *http.Request) {
	userId, err := h.auth.GetUserIDFromAuthenticationToken(r.Cookies())
	if err != nil {
		// TODO: setup proper error page.
		h.ErrorLog.Println(err)
		return
	}

	user, err := h.db.FindUserByID(userId)
	if err != nil {
		// TODO: setup proper error page.
		h.ErrorLog.Println(err)
		return
	}

	fmt.Println(user)

	h.renderers.Page.RenderHTML(w, "profile.page.tmpl", html.ProfilePageData{
		Email: "hello@me.com",
	})
}
