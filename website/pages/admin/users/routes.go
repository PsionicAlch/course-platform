package users

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.UsersGet)
	router.Get("/htmx", handlers.UsersPaginationGet)

	return router
}
