package settings

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.SettingsGet)

	// TODO: Set up the logic for managing a user's whitelisted IP addresses.

	return router
}
