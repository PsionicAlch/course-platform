package rss

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.RSSGet)

	router.Get("/tutorials", handlers.RSSTutorialsGet)
	router.Get("/tutorials/{tutorial-slug}", handlers.RSSTutorialGet)

	router.Get("/courses", handlers.RSSCoursesGet)

	router.Get("/author/{author-slug}/tutorials", handlers.RSSTutorialAuthorGet)
	router.Get("/author/{author-slug}/courses", handlers.RSSCourseAuthorGet)

	return router
}
