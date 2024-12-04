package tutorials

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.TutorialsGet)
	router.Get("/htmx", handlers.TutorialsPaginationGet)

	router.Get("/{slug}", handlers.TutorialGet)
	router.With(handlers.Auth.AllowAuthenticated("")).Post("/{slug}/like", handlers.LikeTutorialPost)
	router.With(handlers.Auth.AllowAuthenticated("")).Post("/{slug}/bookmark", handlers.BookmarkTutorialPost)

	router.Get("/{slug}/comments", handlers.CommentsGet)
	router.With(handlers.Auth.AllowAuthenticated("")).Post("/{slug}/comments", handlers.CommentsPost)

	return router
}
