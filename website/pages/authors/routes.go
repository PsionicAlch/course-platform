package authors

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/website/pages/authors/courses"
	"github.com/PsionicAlch/psionicalch-home/website/pages/authors/tutorials"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	tutorialHandlers := tutorials.SetupHandlers(handlers.Render.Page, handlers.Render.Htmx, handlers.Database)
	courseHandlers := courses.SetupHandlers(handlers.Render.Page, handlers.Render.Htmx, handlers.Database)

	router := chi.NewRouter()

	router.Route("/{author-slug}", func(r chi.Router) {
		r.Get("/", handlers.AuthorGet)

		r.Mount("/tutorials", tutorials.RegisterRoutes(tutorialHandlers))
		r.Mount("/courses", courses.RegisterRoutes(courseHandlers))
	})

	return router
}
