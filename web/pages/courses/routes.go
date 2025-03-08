package courses

import (
	"net/http"

	"github.com/PsionicAlch/course-platform/web/pages"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Use(handlerContext.Authentication.SetUser)
	router.Use(handlerContext.Session.SessionMiddleware)

	router.Get("/", handlers.CoursesGet)
	router.Get("/htmx", handlers.CoursesPaginationGet)

	router.Get("/{course-slug}", handlers.CourseGet)

	router.With(handlerContext.Authentication.AllowAuthenticated("/accounts/login")).Get("/{course-slug}/purchase", handlers.PurchaseCourseGet)
	router.With(handlerContext.Authentication.AllowAuthenticated("/accounts/login")).Post("/{course-slug}/purchase", handlers.PurchaseCoursePost)

	router.Get("/{course-slug}/purchase/success", handlers.PurchaseCourseSuccessGet)
	router.Get("/{course-slug}/purchase/cancel", handlers.PurchaseCourseCancelGet)
	router.Get("/{course-slug}/purchase/check", handlers.PurchaseCourseCheckGet)

	router.With(handlerContext.Authentication.AllowAuthenticated("/accounts/login")).Post("/{course-slug}/purchase/validate", handlers.ValidatePurchasePost)

	return router
}
