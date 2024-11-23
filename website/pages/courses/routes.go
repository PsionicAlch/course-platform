package courses

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.CoursesGet)
	router.Get("/{slug}", handlers.CourseGet)
	router.With(handlers.Auth.AllowAuthenticated("/accounts/login")).Get("/{slug}/purchase", handlers.PurchaseCourseGet)

	return router
}
