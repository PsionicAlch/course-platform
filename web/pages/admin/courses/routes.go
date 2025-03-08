package courses

import (
	"net/http"

	"github.com/PsionicAlch/course-platform/web/pages"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Get("/", handlers.CoursesGet)
	router.Get("/htmx", handlers.CoursesPaginationGet)

	router.Get("/htmx/change-published/{course-id}", handlers.PublishedEditGet)
	router.Post("/htmx/change-published/{course-id}", handlers.PublishedEditPost)

	router.Get("/htmx/change-author/{course-id}", handlers.AuthorEditGet)
	router.Post("/htmx/change-author/{course-id}", handlers.AuthorEditPost)

	return router
}
