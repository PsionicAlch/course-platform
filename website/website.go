package website

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	gocache "github.com/PsionicAlch/psionicalch-home/internal/cache/go-cache"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database"
	"github.com/PsionicAlch/psionicalch-home/internal/payments"
	vanillahtml "github.com/PsionicAlch/psionicalch-home/internal/render/renderers/vanilla-html"
	vanillatext "github.com/PsionicAlch/psionicalch-home/internal/render/renderers/vanilla-text"
	"github.com/PsionicAlch/psionicalch-home/internal/session"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/pkg/sitemapper"
	"github.com/PsionicAlch/psionicalch-home/website/assets"
	"github.com/PsionicAlch/psionicalch-home/website/config"
	"github.com/PsionicAlch/psionicalch-home/website/emails"
	"github.com/PsionicAlch/psionicalch-home/website/generators"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages/accounts"
	"github.com/PsionicAlch/psionicalch-home/website/pages/admin"
	"github.com/PsionicAlch/psionicalch-home/website/pages/authors"
	"github.com/PsionicAlch/psionicalch-home/website/pages/certificates"
	"github.com/PsionicAlch/psionicalch-home/website/pages/courses"
	"github.com/PsionicAlch/psionicalch-home/website/pages/general"
	"github.com/PsionicAlch/psionicalch-home/website/pages/profile"
	"github.com/PsionicAlch/psionicalch-home/website/pages/rss"
	"github.com/PsionicAlch/psionicalch-home/website/pages/settings"
	"github.com/PsionicAlch/psionicalch-home/website/pages/sitemap"
	"github.com/PsionicAlch/psionicalch-home/website/pages/tutorials"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func StartWebsite() {
	// Set up loggers for main.
	loggers := utils.CreateLoggers("WEBSITE")

	// Set up project config.
	err := config.SetupConfig()
	if err != nil {
		loggers.ErrorLog.Fatalln(err)
	}

	// Set up database.
	db, err := sqlite_database.CreateSQLiteDatabase("/db/db.sqlite", "/db/migrations")
	if err != nil {
		loggers.ErrorLog.Fatalln(err)
	}
	defer db.Close()

	// Set up notifications.
	cookieName := config.GetWithoutError[string]("NOTIFICATION_COOKIE_NAME")
	domainName := config.GetWithoutError[string]("DOMAIN_NAME")
	sessions := session.SetupSession(cookieName, domainName)

	// Set up renderers.
	pagesRenderer, err := vanillahtml.SetupVanillaHTMLRenderer(sessions, html.HTMLFiles, ".page.tmpl", "pages", "layouts/*.layout.tmpl", "components/*.component.tmpl")
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to set up pages renderer: ", err)
	}

	htmxRenderer, err := vanillahtml.SetupVanillaHTMLRenderer(nil, html.HTMLFiles, ".htmx.tmpl", "htmx", "components/*.component.tmpl")
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to set up htmx renderer: ", err)
	}

	xmlRenderer, err := vanillatext.SetupVanillaTextRenderer(html.XMLFiles, ".xml.tmpl", "xml")
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to set up rss renderer: ", err)
	}

	emailRenderer, err := vanillahtml.SetupVanillaHTMLRenderer(nil, html.HTMLFiles, ".email.tmpl", "emails", "layouts/email.layout.tmpl")
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to set up email renderer: ", err)
	}

	// Set up authentication system.
	authLifetime := time.Duration(config.GetWithoutError[int]("AUTH_TOKEN_LIFETIME")) * time.Minute
	pwdResetLifetime := time.Minute * 30
	cookieName = config.GetWithoutError[string]("AUTH_COOKIE_NAME")
	currentKey := config.GetWithoutError[string]("CURRENT_SECURE_COOKIE_KEY")
	previousKey := config.GetWithoutError[string]("PREVIOUS_SECURE_COOKIE_KEY")
	auth, err := authentication.SetupAuthentication(db, sessions, authLifetime, pwdResetLifetime, cookieName, domainName, currentKey, previousKey)
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to set up authentication: ", err)
	}

	// Set up emails.
	emailer := emails.SetupEmails(emailRenderer)

	// Set up payments.
	payment := payments.SetupPayments(config.GetWithoutError[string]("STRIPE_SECRET_KEY"), config.GetWithoutError[string]("STRIPE_WEBHOOK_SECRET"), db)

	// Set up cache.
	gens := &gocache.Generators{
		RSSFeed:                generators.RSSFeed(utils.CreateLoggers("RSS FEED GENERATOR"), db, xmlRenderer),
		TutorialsRSSFeed:       generators.TutorialsRSSFeed(utils.CreateLoggers("TUTORIALS RSS FEED GENERATOR"), db, xmlRenderer),
		CoursesRSSFeed:         generators.CoursesRSSFeed(utils.CreateLoggers("COURSES RSS FEED GENERATOR"), db, xmlRenderer),
		AuthorTutorialsRSSFeed: generators.AuthorTutorialsRSSFeed(utils.CreateLoggers("AUTHOR TUTORIALS RSS FEED GENERATOR"), db, xmlRenderer),
		AuthorCoursesRSSFeed:   generators.AuthorCoursesRSSFeed(utils.CreateLoggers("AUTHOR COURSES RSS FEED GENERATOR"), db, xmlRenderer),
	}
	cache := gocache.SetupGoCache(gens)

	// Set up sitemapper.
	mapper := sitemapper.NewSiteMapper("http://localhost:8080", time.Hour*24*7)

	// Set up handlers.
	generalHandlers := general.SetupHandlers(pagesRenderer, db)
	tutorialHandlers := tutorials.SetupHandlers(pagesRenderer, htmxRenderer, db, sessions, auth)
	courseHandlers := courses.SetupHandlers(pagesRenderer, htmxRenderer, db, sessions, auth, payment)
	accountHandlers := accounts.SetupHandlers(pagesRenderer, htmxRenderer, auth, emailer, sessions)
	profileHandlers := profile.SetupHandlers(pagesRenderer, htmxRenderer, auth, db, sessions)
	settingsHandlers := settings.SetupHandlers(pagesRenderer, htmxRenderer, db, sessions, auth)
	adminHandlers := admin.SetupHandlers(pagesRenderer, htmxRenderer, db, auth)
	authorsHandlers := authors.SetupHandlers(pagesRenderer, htmxRenderer, db)
	certificateHandlers := certificates.SetupHandlers(pagesRenderer, db)
	rssHandlers := rss.SetupHandlers(cache)
	sitemapHandlers := sitemap.SetupHandlers(mapper)

	// Create new router.
	router := chi.NewRouter()

	// Set up middleware.
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Register payments webhook.
	router.Post("/payments/webhook", payment.Webhook)

	// Set up asset routes.
	router.Mount("/assets", http.StripPrefix("/assets", http.FileServerFS(assets.AssetFiles)))

	// Set up RSS routes.
	router.Mount("/rss", rss.RegisterRoutes(rssHandlers))
	router.Mount("/sitemap", sitemap.RegisterRoutes(sitemapHandlers))

	// Set up routes.
	router.With(auth.SetUser, sessions.SessionMiddleware).Mount("/", general.RegisterRoutes(generalHandlers))
	router.With(auth.SetUser, sessions.SessionMiddleware).Mount("/accounts", accounts.RegisterRoutes(accountHandlers))
	router.With(auth.SetUser, sessions.SessionMiddleware).Mount("/profile", profile.RegisterRoutes(profileHandlers))
	router.With(auth.SetUser, sessions.SessionMiddleware).Mount("/tutorials", tutorials.RegisterRoutes(tutorialHandlers))
	router.With(auth.SetUser, sessions.SessionMiddleware).Mount("/courses", courses.RegisterRoutes(courseHandlers))
	router.With(auth.SetUserWithEmail(emailer), sessions.SessionMiddleware).Mount("/settings", settings.RegisterRoutes(settingsHandlers))
	router.With(auth.SetUserWithEmail(emailer), sessions.SessionMiddleware).Mount("/admin", admin.RegisterRoutes(adminHandlers))
	router.With(auth.SetUser, sessions.SessionMiddleware).Mount("/authors", authors.RegisterRoutes(authorsHandlers))
	router.With(auth.SetUser, sessions.SessionMiddleware).Mount("/certificates", certificates.RegisterRoutes(certificateHandlers))

	// Set up 404 handler.
	router.With(auth.SetUser, sessions.SessionMiddleware).NotFound(func(w http.ResponseWriter, r *http.Request) {
		user := authentication.GetUserFromRequest(r)

		if err := pagesRenderer.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user)}, http.StatusNotFound); err != nil {
			loggers.ErrorLog.Println(err)
		}
	})

	// Start server.
	port := config.GetWithoutError[string]("PORT")
	loggers.InfoLog.Println("Starting server on port:", port)
	loggers.ErrorLog.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
