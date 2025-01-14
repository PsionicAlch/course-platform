package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	pm "github.com/PsionicAlch/psionicalch-home/internal/middleware"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/web/config"
	"github.com/PsionicAlch/psionicalch-home/web/html"
	"github.com/PsionicAlch/psionicalch-home/web/pages"
	"github.com/PsionicAlch/psionicalch-home/web/pages/accounts"
	"github.com/PsionicAlch/psionicalch-home/web/pages/admin"
	"github.com/PsionicAlch/psionicalch-home/web/pages/authors"
	"github.com/PsionicAlch/psionicalch-home/web/pages/certificates"
	"github.com/PsionicAlch/psionicalch-home/web/pages/courses"
	"github.com/PsionicAlch/psionicalch-home/web/pages/general"
	"github.com/PsionicAlch/psionicalch-home/web/pages/profile"
	"github.com/PsionicAlch/psionicalch-home/web/pages/rss"
	"github.com/PsionicAlch/psionicalch-home/web/pages/settings"
	"github.com/PsionicAlch/psionicalch-home/web/pages/sitemap"
	"github.com/PsionicAlch/psionicalch-home/web/pages/tutorials"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
)

func StartWeb() {
	// Set up loggers for main.
	loggers := utils.CreateLoggers("WEBSITE")

	handlerContext, err := pages.CreateHandlerContext()
	if err != nil {
		loggers.ErrorLog.Fatalln(err)
	}

	// Create new router.
	router := chi.NewRouter()

	// Set up middleware.
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(pm.CSRFProtection)
	router.Use(pm.RateLimiter(50, time.Minute, handlerContext.Renderers.Page))

	// Register payments webhook.
	router.Post("/payments/webhook", handlerContext.Payment.Webhook)

	// Set up routes.
	router.Mount("/rss", rss.RegisterRoutes(handlerContext))
	router.Mount("/sitemap", sitemap.RegisterRoutes(handlerContext))
	router.Mount("/", general.RegisterRoutes(handlerContext))
	router.Mount("/accounts", accounts.RegisterRoutes(handlerContext))
	router.Mount("/profile", profile.RegisterRoutes(handlerContext))
	router.Mount("/tutorials", tutorials.RegisterRoutes(handlerContext))
	router.Mount("/courses", courses.RegisterRoutes(handlerContext))
	router.Mount("/settings", settings.RegisterRoutes(handlerContext))
	router.Mount("/admin", admin.RegisterRoutes(handlerContext))
	router.Mount("/authors", authors.RegisterRoutes(handlerContext))
	router.Mount("/certificates", certificates.RegisterRoutes(handlerContext))

	// Set up 404 handler.
	router.With(handlerContext.Authentication.SetUser, handlerContext.Session.SessionMiddleware).NotFound(func(w http.ResponseWriter, r *http.Request) {
		user := authentication.GetUserFromRequest(r)

		if err := handlerContext.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusNotFound); err != nil {
			loggers.ErrorLog.Println(err)
		}
	})

	// Start server.
	port := config.GetWithoutError[string]("PORT")
	loggers.InfoLog.Println("Starting server on port:", port)
	loggers.ErrorLog.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
