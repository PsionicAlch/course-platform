package courses

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.CoursesGet)
	router.Get("/htmx", handlers.CoursesPaginationGet)

	router.Get("/{slug}", handlers.CourseGet)
	router.Get("/{course_slug}/{chapter_slug}", handlers.CourseChapterGet)

	return router
}
