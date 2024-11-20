package tutorials

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.TutorialsGet)
	router.Get("/page/{page-number}", handlers.TutorialsPaginationGet)
	router.Get("/search", handlers.TutorialsSearchGet)

	router.Get("/{slug}", handlers.TutorialGet)
	router.With(handlers.Auth.AllowAuthenticated("")).Post("/{slug}/like", handlers.LikeTutorialPost)
	router.With(handlers.Auth.AllowAuthenticated("")).Post("/{slug}/bookmark", handlers.BookmarkTutorialPost)

	return router
}
