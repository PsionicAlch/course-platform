package general

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.HomeGet)
	router.Get("/rss", handlers.RSSGet)
	router.Get("/affiliate-program", handlers.AffiliateProgramGet)
	router.Get("/privacy-policy", handlers.PrivacyPolicyGet)
	router.Get("/refund-policy", handlers.RefundPolicyGet)

	return router
}
