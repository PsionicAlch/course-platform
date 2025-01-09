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
		r.Use(handlers.UserBoughtCourse)

		r.Get("/", handlers.CourseGet)
		r.Get("/certificate", handlers.CourseCertificateGet)
		r.Get("/{chapter-slug}", handlers.CourseChapterGet)
		r.Post("/{chapter-slug}/finish", handlers.CourseChapterFinishPost)
	})

	return router
}
