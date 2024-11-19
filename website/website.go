package website

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database"
	"github.com/PsionicAlch/psionicalch-home/internal/render/renderers/vanilla"
	"github.com/PsionicAlch/psionicalch-home/internal/session"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/assets"
	"github.com/PsionicAlch/psionicalch-home/website/config"
	"github.com/PsionicAlch/psionicalch-home/website/emails"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages/accounts"
	"github.com/PsionicAlch/psionicalch-home/website/pages/admin"
	"github.com/PsionicAlch/psionicalch-home/website/pages/authors"
	"github.com/PsionicAlch/psionicalch-home/website/pages/courses"
	"github.com/PsionicAlch/psionicalch-home/website/pages/general"
	"github.com/PsionicAlch/psionicalch-home/website/pages/profile"
	"github.com/PsionicAlch/psionicalch-home/website/pages/settings"
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
	pagesRenderer, err := vanilla.SetupVanillaRenderer(sessions, html.HTMLFiles, ".page.tmpl", "pages", "layouts/*.layout.tmpl", "components/*.component.tmpl")
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to set up pages renderer: ", err)
	}

	htmxRenderer, err := vanilla.SetupVanillaRenderer(nil, html.HTMLFiles, ".htmx.tmpl", "htmx", "components/*.component.tmpl")
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to set up htmx renderer: ", err)
	}

	emailRenderer, err := vanilla.SetupVanillaRenderer(nil, html.HTMLFiles, ".email.tmpl", "emails", "layouts/email.layout.tmpl")
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

	// Set up handlers.
	generalHandlers := general.SetupHandlers(pagesRenderer, db)
	tutorialHandlers := tutorials.SetupHandlers(pagesRenderer, htmxRenderer, db, sessions)
	courseHandlers := courses.SetupHandlers(pagesRenderer)
	accountHandlers := accounts.SetupHandlers(pagesRenderer, htmxRenderer, auth, emailer, sessions)
	profileHandlers := profile.SetupHandlers(pagesRenderer, auth, db)
	settingsHandlers := settings.SetupHandlers(pagesRenderer)
	adminHandlers := admin.SetupHandlers(pagesRenderer, auth)
	authorsHandlers := authors.SetupHandlers(pagesRenderer)

	// Create new router.
	router := chi.NewRouter()

	// Set up middleware.
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set up 404 handler.
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Set up 404 page.
		pagesRenderer.RenderHTML(w, r.Context(), "404.page.tmpl", nil, http.StatusNotFound)
	})

	// Set up asset routes.
	router.Mount("/assets", http.StripPrefix("/assets", http.FileServerFS(assets.AssetFiles)))

	// Set up routes.
	router.With(auth.SetUser, sessions.SessionMiddleware).Mount("/", general.RegisterRoutes(generalHandlers))
	router.With(auth.SetUser, sessions.SessionMiddleware).Mount("/accounts", accounts.RegisterRoutes(accountHandlers))
	router.With(auth.SetUser, sessions.SessionMiddleware).Mount("/profile", profile.RegisterRoutes(profileHandlers))
	router.With(auth.SetUser, sessions.SessionMiddleware).Mount("/tutorials", tutorials.RegisterRoutes(tutorialHandlers))
	router.With(auth.SetUser, sessions.SessionMiddleware).Mount("/courses", courses.RegisterRoutes(courseHandlers))
	router.With(auth.SetUserWithEmail(emailer), sessions.SessionMiddleware).Mount("/settings", settings.RegisterRoutes(settingsHandlers))
	router.With(auth.SetUserWithEmail(emailer), sessions.SessionMiddleware).Mount("/admin", admin.RegisterRoutes(adminHandlers))
	router.With(auth.SetUser, sessions.SessionMiddleware).Mount("/authors", authors.RegisterRoutes(authorsHandlers))

	// Start server.
	port := config.GetWithoutError[string]("PORT")
	loggers.InfoLog.Println("Starting server on port:", port)
	loggers.ErrorLog.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
