package admin

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/web/pages"
	"github.com/PsionicAlch/psionicalch-home/web/pages/admin/comments"
	"github.com/PsionicAlch/psionicalch-home/web/pages/admin/courses"
	"github.com/PsionicAlch/psionicalch-home/web/pages/admin/discounts"
	"github.com/PsionicAlch/psionicalch-home/web/pages/admin/purchases"
	"github.com/PsionicAlch/psionicalch-home/web/pages/admin/refunds"
	"github.com/PsionicAlch/psionicalch-home/web/pages/admin/tutorials"
	"github.com/PsionicAlch/psionicalch-home/web/pages/admin/users"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Use(handlerContext.Authentication.SetUserWithEmail(handlerContext.Emailer))
	router.Use(handlerContext.Session.SessionMiddleware)
	router.Use(handlerContext.Authentication.AllowAdmin("/"))

	router.Get("/", handlers.AdminGet)

	router.Mount("/comments", comments.RegisterRoutes(handlerContext))
	router.Mount("/courses", courses.RegisterRoutes(handlerContext))
	router.Mount("/discounts", discounts.RegisterRoutes(handlerContext))
	router.Mount("/purchases", purchases.RegisterRoutes(handlerContext))
	router.Mount("/refunds", refunds.RegisterRoutes(handlerContext))
	router.Mount("/tutorials", tutorials.RegisterRoutes(handlerContext))
	router.Mount("/users", users.RegisterRoutes(handlerContext))

	return router
}
