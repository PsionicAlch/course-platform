package comments

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.CommentsGet)
	router.Get("/htmx", handlers.CommentsPaginationGet)
	router.Delete("/{comment-id}", handlers.CommentDelete)

	return router
}
