package website

import (
	"fmt"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/config"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database"
	"github.com/PsionicAlch/psionicalch-home/internal/render/renderers/vanilla"
	scssession "github.com/PsionicAlch/psionicalch-home/internal/session/scs_session"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/pkg/gatekeeper"
	"github.com/PsionicAlch/psionicalch-home/website/assets"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages/accounts"
	"github.com/PsionicAlch/psionicalch-home/website/pages/courses"
	"github.com/PsionicAlch/psionicalch-home/website/pages/general"
	"github.com/PsionicAlch/psionicalch-home/website/pages/profile"
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

	// Set up session.
	session := scssession.NewSession()

	// Set up renderers.
	pagesRenderer, err := vanilla.SetupVanillaRenderer(html.HTMLFiles, "pages", "layouts/*.layout.tmpl", "components/*.component.tmpl")
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to set up pages renderer: ", err)
	}

	htmxRenderer, err := vanilla.SetupVanillaRenderer(html.HTMLFiles, "htmx", "components/*.component.tmpl")
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to set up htmx renderer: ", err)
	}

	// Set up tutorials.
	tuts := tutorials.SetupTutorials()

	// Set up Gatekeeper.
	authCookieName := config.GetWithoutError[string]("AUTH_COOKIE_NAME")
	websiteDomain := config.GetWithoutError[string]("DOMAIN_NAME")
	authLifetime := config.GetWithoutError[int]("AUTH_TOKEN_LIFETIME")
	currentGatekeeperKey := config.GetWithoutError[string]("GATEKEEPER_CURRENT_KEY")
	prevGatekeeperKey := config.GetWithoutError[string]("GATEKEEPER_PREVIOUS_KEY")

	auth, err := gatekeeper.NewGatekeeper(authCookieName, websiteDomain, authLifetime, currentGatekeeperKey, prevGatekeeperKey, db)
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to set up authentication: ", err)
	}

	// Set up handlers.
	generalHandlers := general.SetupHandlers(pagesRenderer)
	accountHandlers := accounts.SetupHandlers(pagesRenderer, htmxRenderer, session, auth)
	tutorialHandlers := tutorials.SetupHandlers(pagesRenderer, tuts)
	courseHandlers := courses.SetupHandlers(pagesRenderer)
	profileHandlers := profile.SetupHandlers(pagesRenderer, auth)

	// Create new router.
	router := chi.NewRouter()

	// Set up middleware.
	router.Use(middleware.RealIP)
	router.Use(session.LoadSession)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set up 404 handler.
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		pagesRenderer.RenderHTML(w, "404.page.tmpl", nil, http.StatusNotFound)
	})

	// Set up asset routes.
	router.Mount("/assets", http.StripPrefix("/assets", http.FileServerFS(assets.AssetFiles)))

	// Set up routes.
	router.Mount("/", general.RegisterRoutes(generalHandlers))
	router.Mount("/accounts", accounts.RegisterRoutes(accountHandlers))
	router.Mount("/profile", profile.RegisterRoutes(profileHandlers))
	router.Mount("/tutorials", tutorials.RegisterRoutes(tutorialHandlers))
	router.Mount("/courses", courses.RegisterRoutes(courseHandlers))

	// Start server.
	port := config.GetWithoutError[string]("PORT")
	loggers.InfoLog.Println("Starting server on port:", port)
	loggers.ErrorLog.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
