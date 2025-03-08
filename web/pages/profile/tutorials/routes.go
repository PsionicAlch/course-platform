package tutorials

import (
	"net/http"

	"github.com/PsionicAlch/course-platform/web/pages"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Get("/bookmarks", handlers.TutorialsBookmarksGet)
	router.Get("/bookmarks/htmx", handlers.TutorialsBookmarksPaginationGet)

	router.Get("/likes", handlers.TutorialsLikesGet)
	router.Get("/likes/htmx", handlers.TutorialsLikesPaginationGet)

	return router
}
