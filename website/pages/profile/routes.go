package profile

import (
	"net/http"

	affiliatehistory "github.com/PsionicAlch/psionicalch-home/website/pages/profile/affiliate-history"
	"github.com/PsionicAlch/psionicalch-home/website/pages/profile/courses"
	"github.com/PsionicAlch/psionicalch-home/website/pages/profile/tutorials"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	affiliateHistoryHandlers := affiliatehistory.SetupHandlers(handlers.Renderers.Page, handlers.Renderers.Htmx, handlers.Database)
	courseHandlers := courses.SetupHandlers(handlers.Renderers.Page, handlers.Renderers.Htmx, handlers.Database)
	tutorialHandlers := tutorials.SetupHandlers(handlers.Renderers.Page, handlers.Renderers.Htmx, handlers.Database)

	router := chi.NewRouter()

	router.Use(handlers.Auth.AllowAuthenticated("/accounts/login"))

	router.Get("/", handlers.ProfileGet)

	router.Mount("/affiliate-history", affiliatehistory.RegisterRoutes(affiliateHistoryHandlers))
	router.Mount("/courses", courses.RegisterRoutes(courseHandlers))
	router.Mount("/tutorials", tutorials.RegisterRoutes(tutorialHandlers))

	return router
}
