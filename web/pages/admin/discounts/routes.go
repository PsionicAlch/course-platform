package discounts

import (
	"net/http"

	"github.com/PsionicAlch/course-platform/web/pages"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Get("/", handlers.DiscountsGet)
	router.Get("/htmx", handlers.DiscountsPaginationGet)

	router.Post("/add", handlers.NewDiscountPost)
	router.Post("/validate/add", handlers.ValidateNewDiscountPost)
	router.Get("/validate/empty", handlers.EmptyNewDiscountGet)

	router.Get("/htmx/change-status/{discount-id}", handlers.StatusEditGet)
	router.Post("/htmx/change-status/{discount-id}", handlers.StatusEditPost)

	return router
}
