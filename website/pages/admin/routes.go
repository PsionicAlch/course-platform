package admin

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/website/pages/admin/authors"
	"github.com/PsionicAlch/psionicalch-home/website/pages/admin/comments"
	"github.com/PsionicAlch/psionicalch-home/website/pages/admin/courses"
	"github.com/PsionicAlch/psionicalch-home/website/pages/admin/discounts"
	"github.com/PsionicAlch/psionicalch-home/website/pages/admin/purchases"
	"github.com/PsionicAlch/psionicalch-home/website/pages/admin/refunds"
	"github.com/PsionicAlch/psionicalch-home/website/pages/admin/tutorials"
	"github.com/PsionicAlch/psionicalch-home/website/pages/admin/users"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	authorsHandlers := authors.SetupHandlers(handlers.Renderers.Page, handlers.Renderers.Htmx, handlers.Database, handlers.Auth)
	commentsHandlers := comments.SetupHandlers(handlers.Renderers.Page, handlers.Renderers.Htmx, handlers.Database, handlers.Auth)
	coursesHandlers := courses.SetupHandlers(handlers.Renderers.Page, handlers.Renderers.Htmx, handlers.Database, handlers.Auth)
	discountsHandlers := discounts.SetupHandlers(handlers.Renderers.Page, handlers.Renderers.Htmx, handlers.Database, handlers.Auth)
	purchasesHandlers := purchases.SetupHandlers(handlers.Renderers.Page, handlers.Renderers.Htmx, handlers.Database, handlers.Auth)
	refundsHandlers := refunds.SetupHandlers(handlers.Renderers.Page, handlers.Renderers.Htmx, handlers.Database, handlers.Auth)
	tutorialsHandlers := tutorials.SetupHandlers(handlers.Renderers.Page, handlers.Renderers.Htmx, handlers.Database, handlers.Auth)
	usersHandlers := users.SetupHandlers(handlers.Renderers.Page, handlers.Renderers.Htmx, handlers.Database, handlers.Auth)

	router := chi.NewRouter()

	router.Use(handlers.Auth.AllowAdmin("/"))

	router.Get("/", handlers.AdminGet)

	router.Mount("/authors", authors.RegisterRoutes(authorsHandlers))
	router.Mount("/comments", comments.RegisterRoutes(commentsHandlers))
	router.Mount("/courses", courses.RegisterRoutes(coursesHandlers))
	router.Mount("/discounts", discounts.RegisterRoutes(discountsHandlers))
	router.Mount("/purchases", purchases.RegisterRoutes(purchasesHandlers))
	router.Mount("/refunds", refunds.RegisterRoutes(refundsHandlers))
	router.Mount("/tutorials", tutorials.RegisterRoutes(tutorialsHandlers))
	router.Mount("/users", users.RegisterRoutes(usersHandlers))

	return router
}
