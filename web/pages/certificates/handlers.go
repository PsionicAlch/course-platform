package certificates

import (
	"net/http"

	"github.com/PsionicAlch/course-platform/internal/authentication"
	"github.com/PsionicAlch/course-platform/internal/database"
	"github.com/PsionicAlch/course-platform/internal/utils"
	"github.com/PsionicAlch/course-platform/web/html"
	"github.com/PsionicAlch/course-platform/web/pages"
	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
)

type Handlers struct {
	utils.Loggers
	*pages.HandlerContext
}

func SetupHandlers(handlerContext *pages.HandlerContext) *Handlers {
	loggers := utils.CreateLoggers("CERTIFICATE HANDLERS")

	return &Handlers{
		Loggers:        loggers,
		HandlerContext: handlerContext,
	}
}

func (h *Handlers) CertificateGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := &html.CertificatePage{}

	certificateId := chi.URLParam(r, "certificate-id")

	certificate, err := h.Database.GetCertificateFromID(certificateId)
	if err != nil {
		h.ErrorLog.Printf("Failed to get certificate (\"%s\"): %s\n", certificateId, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if certificate == nil {
		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusNotFound); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Certificate = certificate

	certificateUser, err := h.Database.GetUserByID(certificate.UserID, database.All)
	if err != nil {
		h.ErrorLog.Printf("Failed to get user from certificate (\"%s\"): %s\n", certificate.UserID, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if certificateUser == nil {
		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusNotFound); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.User = certificateUser

	course, err := h.Database.GetCourseByID(certificate.CourseID)
	if err != nil {
		h.ErrorLog.Printf("Failed to get course from certificate (\"%s\"): %s\n", certificate.CourseID, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if course == nil {
		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusNotFound); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Course = course

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "certificate", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}
