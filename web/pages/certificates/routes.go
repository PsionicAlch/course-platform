package certificates

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/web/pages"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Use(handlerContext.Authentication.SetUser)
	router.Use(handlerContext.Session.SessionMiddleware)

	router.Get("/{certificate-id}", handlers.CertificateGet)

	return router
}
