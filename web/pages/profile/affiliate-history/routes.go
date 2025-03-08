package affiliatehistory

import (
	"net/http"

	"github.com/PsionicAlch/course-platform/web/pages"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Get("/", handlers.AffiliateHistoryGet)
	router.Get("/htmx", handlers.AffiliateHistoryPaginationGet)

	return router
}
