package website

import (
	"fmt"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/config"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database"
	"github.com/PsionicAlch/psionicalch-home/internal/render/renderers/vanilla"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/assets"
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

	// Set up renderers.
	pagesRenderer, err := vanilla.SetupVanillaRenderer(html.HTMLFiles, "pages", "layouts/*.layout.tmpl", "components/*.component.tmpl")
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to set up pages renderer: ", err)
	}

	htmxRenderer, err := vanilla.SetupVanillaRenderer(html.HTMLFiles, "htmx", "components/*.component.tmpl")
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to set up htmx renderer: ", err)
	}

	// Set up authentication system.
	auth, err := authentication.SetupAuthentication(db)
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to set up authentication: ", err)
	}

	// Set up handlers.
	generalHandlers := general.SetupHandlers(pagesRenderer, db)
	tutorialHandlers := tutorials.SetupHandlers(pagesRenderer, db)
	courseHandlers := courses.SetupHandlers(pagesRenderer)
	accountHandlers := accounts.SetupHandlers(pagesRenderer, htmxRenderer, auth)
	profileHandlers := profile.SetupHandlers(pagesRenderer, db)
	settingsHandlers := settings.SetupHandlers(pagesRenderer)
	adminHandlers := admin.SetupHandlers(pagesRenderer)
	authorsHandlers := authors.SetupHandlers(pagesRenderer)

	// Create new router.
	router := chi.NewRouter()

	// Set up middleware.
	router.Use(middleware.RealIP)
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
	router.Mount("/settings", settings.RegisterRoutes(settingsHandlers))
	router.Mount("/admin", admin.RegisterRoutes(adminHandlers))
	router.Mount("/authors", authors.RegisterRoutes(authorsHandlers))

	// Start server.
	port := config.GetWithoutError[string]("PORT")
	loggers.InfoLog.Println("Starting server on port:", port)
	loggers.ErrorLog.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
