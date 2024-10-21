package main

import (
	"fmt"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/config"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database"
	"github.com/PsionicAlch/psionicalch-home/internal/services/accounts"
	"github.com/PsionicAlch/psionicalch-home/internal/services/assets"
	"github.com/PsionicAlch/psionicalch-home/internal/services/courses"
	"github.com/PsionicAlch/psionicalch-home/internal/services/general"
	"github.com/PsionicAlch/psionicalch-home/internal/services/tutorials"
	scssession "github.com/PsionicAlch/psionicalch-home/internal/session/scs_session"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/internal/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Set up loggers for main.
	loggers := utils.CreateLoggers("MAIN")

	// Set up project config.
	config.NewConfig()

	// Set up database.
	db := sqlite_database.CreateSQLiteDatabase()
	defer db.Close()

	// Set up session.
	session := scssession.NewSession()

	// Set up views.
	view := views.SetupRenderer()

	// Set up tutorials.
	tuts := tutorials.SetupTutorials()

	// Set up authentication.
	auth := authentication.CreateAuthentication(db)

	// Set up handlers.
	generalHandlers := general.SetupHandlers(view)
	accountHandlers := accounts.SetupHandlers(view, session, auth)
	tutorialHandlers := tutorials.SetupHandlers(view, tuts)
	courseHandlers := courses.SetupHandlers(view)

	// Create new router.
	router := chi.NewRouter()

	// Set up middleware.
	router.Use(session.LoadSession)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)

	// Set up 404 handler.
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		view.RenderNotFound(w)
	})

	// Set up routes.
	router.Mount("/", general.RegisterRoutes(generalHandlers))
	router.Mount("/accounts", accounts.RegisterRoutes(accountHandlers))
	router.Mount("/tutorials", tutorials.RegisterRoutes(tutorialHandlers))
	router.Mount("/courses", courses.RegisterRoutes(courseHandlers))
	router.Mount("/assets", assets.RegisterAssetRoutes())

	// Start server.
	loggers.InfoLog.Println("Starting server on port:", config.GetPort())
	loggers.ErrorLog.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.GetPort()), router))
}
