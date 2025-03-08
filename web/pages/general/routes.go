package general

import (
	"net/http"

	"github.com/PsionicAlch/course-platform/web/pages"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Use(handlers.Authentication.SetUser)
	router.Use(handlers.Session.SessionMiddleware)

	router.Get("/", handlers.HomeGet)
	router.Get("/affiliate-program", handlers.AffiliateProgramGet)
	router.Get("/privacy-policy", handlers.PrivacyPolicyGet)
	router.Get("/refund-policy", handlers.RefundPolicyGet)

	return router
}
