package profile

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Use(handlers.auth.AllowAuthenticated("/accounts/login"))

	router.Get("/", handlers.ProfileGet)

	return router
}
