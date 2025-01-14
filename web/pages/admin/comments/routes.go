package comments

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/web/pages"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Get("/", handlers.CommentsGet)
	router.Get("/htmx", handlers.CommentsPaginationGet)
	router.Delete("/{comment-id}", handlers.CommentDelete)

	return router
}
