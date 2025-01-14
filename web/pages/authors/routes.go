package authors

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/web/pages"
	"github.com/PsionicAlch/psionicalch-home/web/pages/authors/courses"
	"github.com/PsionicAlch/psionicalch-home/web/pages/authors/tutorials"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Use(handlerContext.Authentication.SetUser)
	router.Use(handlerContext.Session.SessionMiddleware)

	router.Route("/{author-slug}", func(r chi.Router) {
		r.Get("/", handlers.AuthorGet)

		r.Mount("/tutorials", tutorials.RegisterRoutes(handlerContext))
		r.Mount("/courses", courses.RegisterRoutes(handlerContext))
	})

	return router
}
