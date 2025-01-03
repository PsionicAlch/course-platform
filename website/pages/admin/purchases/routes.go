package purchases

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.PurchasesGet)
	router.Get("/htmx", handlers.PurchasesPaginationGet)

	return router
}
