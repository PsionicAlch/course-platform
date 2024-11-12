package authors

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/{slug}/tutorials", handlers.TutorialsGet)
	router.Get("/{slug}/courses", handlers.CoursesGet)

	return router
}
