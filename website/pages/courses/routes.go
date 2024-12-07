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
	router.With(handlers.Auth.AllowAuthenticated("/accounts/login")).Get("/{course-slug}/purchase", handlers.PurchaseCourseGet)
	router.With(handlers.Auth.AllowAuthenticated("/accounts/login")).Post("/{course-slug}/purchase/", handlers.PurchaseCoursePost)
	router.With(handlers.Auth.AllowAuthenticated("/accounts/login")).Get("/{course-slug}/purchase/success", handlers.PurchaseCourseSuccessGet)
	router.With(handlers.Auth.AllowAuthenticated("/accounts/login")).Get("/{course-slug}/purchase/cancel", handlers.PurchaseCourseCancelGet)
	router.With(handlers.Auth.AllowAuthenticated("/accounts/login")).Post("/{course-slug}/purchase/validate", handlers.ValidatePurchasePost)

	return router
}
