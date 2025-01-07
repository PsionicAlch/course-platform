package general

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
	"github.com/justinas/nosurf"
)

type Handlers struct {
	utils.Loggers
	Renderers *pages.Renderers
	Database  database.Database
}

func SetupHandlers(pageRenderer render.Renderer, db database.Database) *Handlers {
	loggers := utils.CreateLoggers("GENERAL HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: pages.CreateRenderers(pageRenderer, nil, nil),
		Database:  db,
	}
}

func (h *Handlers) HomeGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.GeneralHomePage{
		BasePage: html.NewBasePage(user, nosurf.Token(r)),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "general-home", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) AffiliateProgramGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.GeneralAffiliateProgramPage{
		BasePage: html.NewBasePage(user, nosurf.Token(r)),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "general-affiliate-program", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PrivacyPolicyGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.GeneralPrivacyPolicyPage{
		BasePage: html.NewBasePage(user, nosurf.Token(r)),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "general-privacy-policy", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) RefundPolicyGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.GeneralRefundPolicyPage{
		BasePage: html.NewBasePage(user, nosurf.Token(r)),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "general-refund-policy", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}
