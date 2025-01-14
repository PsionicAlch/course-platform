package profile

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/web/pages"
	affiliatehistory "github.com/PsionicAlch/psionicalch-home/web/pages/profile/affiliate-history"
	"github.com/PsionicAlch/psionicalch-home/web/pages/profile/courses"
	"github.com/PsionicAlch/psionicalch-home/web/pages/profile/tutorials"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Use(handlerContext.Authentication.SetUser)
	router.Use(handlerContext.Session.SessionMiddleware)
	router.Use(handlerContext.Authentication.AllowAuthenticated("/accounts/login"))

	router.Get("/", handlers.ProfileGet)

	router.Mount("/affiliate-history", affiliatehistory.RegisterRoutes(handlerContext))
	router.Mount("/courses", courses.RegisterRoutes(handlerContext))
	router.Mount("/tutorials", tutorials.RegisterRoutes(handlerContext))

	return router
}
