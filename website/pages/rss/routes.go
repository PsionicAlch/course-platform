package rss

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.RSSGet)

	router.Get("/tutorials", handlers.RSSTutorialsGet)
	router.Get("/tutorials/{author-slug}", handlers.RSSTutorialAuthorGet)

	router.Get("/courses", handlers.RSSCoursesGet)
	router.Get("/courses/{author-slug}", handlers.RSSCourseAuthorGet)

	return router
}
