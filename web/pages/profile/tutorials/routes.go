package tutorials

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/bookmarks", handlers.TutorialsBookmarksGet)
	router.Get("/bookmarks/htmx", handlers.TutorialsBookmarksPaginationGet)

	router.Get("/likes", handlers.TutorialsLikesGet)
	router.Get("/likes/htmx", handlers.TutorialsLikesPaginationGet)

	return router
}
