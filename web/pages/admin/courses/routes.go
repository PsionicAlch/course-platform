package courses

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.CoursesGet)
	router.Get("/htmx", handlers.CoursesPaginationGet)

	router.Get("/htmx/change-published/{course-id}", handlers.PublishedEditGet)
	router.Post("/htmx/change-published/{course-id}", handlers.PublishedEditPost)

	router.Get("/htmx/change-author/{course-id}", handlers.AuthorEditGet)
	router.Post("/htmx/change-author/{course-id}", handlers.AuthorEditPost)

	return router
}
