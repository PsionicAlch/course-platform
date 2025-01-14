package tutorials

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/web/pages"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Use(handlerContext.Authentication.SetUser)
	router.Use(handlerContext.Session.SessionMiddleware)

	router.Get("/", handlers.TutorialsGet)
	router.Get("/htmx", handlers.TutorialsPaginationGet)

	router.Get("/{slug}", handlers.TutorialGet)
	router.With(handlerContext.Authentication.AllowAuthenticated("")).Post("/{slug}/like", handlers.LikeTutorialPost)
	router.With(handlerContext.Authentication.AllowAuthenticated("")).Post("/{slug}/bookmark", handlers.BookmarkTutorialPost)

	router.Get("/{slug}/comments", handlers.CommentsGet)
	router.With(handlerContext.Authentication.AllowAuthenticated("")).Post("/{slug}/comments", handlers.CommentsPost)

	return router
}
