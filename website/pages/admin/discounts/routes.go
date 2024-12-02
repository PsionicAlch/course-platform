package discounts

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
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
