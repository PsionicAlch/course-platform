package general

import (
	"net/http"

	"github.com/PsionicAlch/course-platform/internal/authentication"
	"github.com/PsionicAlch/course-platform/internal/utils"
	"github.com/PsionicAlch/course-platform/web/html"
	"github.com/PsionicAlch/course-platform/web/pages"
	"github.com/justinas/nosurf"
)

type Handlers struct {
	utils.Loggers
	*pages.HandlerContext
}

func SetupHandlers(handlerContext *pages.HandlerContext) *Handlers {
	loggers := utils.CreateLoggers("GENERAL HANDLERS")

	return &Handlers{
		Loggers:        loggers,
		HandlerContext: handlerContext,
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
