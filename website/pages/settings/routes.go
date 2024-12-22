package settings

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.SettingsGet)

	router.Route("/validate", func(r chi.Router) {
		r.Post("/change-first-name", handlers.ValidateChangeFirstName)
		r.Post("/change-last-name", handlers.ValidateChangeLastName)
		r.Post("/change-email", handlers.ValidateChangeEmail)
	})

	return router
}
