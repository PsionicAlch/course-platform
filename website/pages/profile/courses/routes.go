package courses

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Base URL: /profile/courses
func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.CoursesGet)
	router.Get("/htmx", handlers.CoursesPaginationGet)

	router.Route("/{course-slug}", func(r chi.Router) {
		router.Use(handlers.UserBoughtCourse)

		router.Get("/", handlers.CourseGet)
		router.Get("/certificate", handlers.CourseCertificateGet)
		router.Get("/{chapter-slug}", handlers.CourseChapterGet)
	})

	return router
}
