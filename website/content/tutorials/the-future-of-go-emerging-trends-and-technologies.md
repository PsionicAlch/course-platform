---
title: "The Future of Go: Emerging Trends and Technologies"
description: "Explore the future of Go and the latest trends in the language."
thumbnail_url: "https://d2ra3i7lkx098s.cloudfront.net/img/golang-course-platform-image.png"
banner_url: "https://d2ra3i7lkx098s.cloudfront.net/img/golang-course-platform-image.png"
keywords: ["keyword 1", "keyword 2", "keyword 3"]

key: "aqeEamtbZMfutfvpymfrjagdjgZPK7P0--YBlB6g6cnvOUDsH_R5Q-c-1XcBAbNlFV1-eYVE8OPPRgyURs732w"
---

## Lorem Ipsum Dolor Sit Amet

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer placerat ante ut sodales venenatis. Donec non aliquam ligula, sit amet sollicitudin ligula. Mauris ornare neque consequat, viverra quam non, fermentum leo. Nullam feugiat ullamcorper ipsum, vel dapibus elit porttitor eget. Cras rutrum, lacus at laoreet convallis, erat libero gravida ex, id accumsan augue magna id nunc. Proin semper elit eu arcu elementum, non aliquet massa scelerisque. Duis porta commodo dapibus. Nulla vestibulum tristique blandit. Suspendisse pulvinar tortor purus, sit amet varius urna vehicula sit amet. Aenean congue mi mi, a blandit lectus maximus efficitur. Nunc quis metus lorem. Morbi et diam in elit laoreet bibendum. Morbi ultrices ligula at massa accumsan, a congue enim ullamcorper. Phasellus sit amet ultricies dolor.

![Alt text for image](https://d2ra3i7lkx098s.cloudfront.net/img/golang-course-platform-image.png "a title for the image")

## Secondary Title

Quisque eget magna sapien. Quisque a nulla at risus convallis dapibus ac pretium eros. Nullam ultrices suscipit orci, non ultrices risus faucibus non. Nullam vulputate laoreet metus, ut consequat urna volutpat ac. Cras dictum dignissim augue in tempus. Duis augue sapien, tincidunt id lectus a, tincidunt sagittis purus. Sed pellentesque elementum ante, non laoreet quam elementum sit amet. Aliquam dictum velit dictum, varius neque bibendum, mattis mi. Cras in dictum purus.

## Secondary Title

Sed ac orci elementum, hendrerit massa non, placerat ligula. Quisque ac lacus mi. Cras nec mi in massa tempus placerat. Nullam gravida justo in diam vehicula, at facilisis ex auctor. Sed sodales sodales feugiat. Quisque eleifend vulputate felis nec tincidunt. Quisque ac pellentesque purus. Nunc ac sodales erat. Etiam venenatis malesuada purus, quis lobortis mauris maximus eu. Fusce elementum, leo in mattis porttitor, erat mauris dignissim odio, at luctus sapien velit quis elit. Quisque sollicitudin dignissim sapien, vitae ultricies neque. Nunc bibendum ut odio ut sagittis. Cras pellentesque felis eu interdum sodales. Curabitur sit amet mi in nisi dapibus tincidunt.

This is what a code block would look like:

```golang
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
	generalHandlers := general.SetupHandlers(pagesRenderer, auth, db)
	accountHandlers := accounts.SetupHandlers(pagesRenderer, htmxRenderer, session, auth)
	tutorialHandlers := tutorials.SetupHandlers(pagesRenderer, auth, db)
	courseHandlers := courses.SetupHandlers(pagesRenderer)
	profileHandlers := profile.SetupHandlers(pagesRenderer, auth, db)

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
```

### Tersary Title

Ut nec viverra lorem. Sed venenatis erat sit amet nisl finibus suscipit at non elit. Nullam vitae ultrices libero. Duis non semper felis. Sed sapien lacus, ullamcorper interdum fermentum ac, consectetur ac felis. Suspendisse ultrices fringilla feugiat. Etiam feugiat vel mi ut efficitur. Sed nisl nulla, pellentesque non dolor non, facilisis accumsan libero. Aliquam porta augue lacus, sit amet porttitor risus tristique tristique. Fusce et pellentesque risus, at aliquam mi. Nam convallis fringilla ligula ac volutpat. Duis ligula magna, ornare ut suscipit ut, dictum non turpis. Fusce dictum urna vel lorem tincidunt, id tempor felis sagittis. Praesent interdum viverra arcu, vitae fringilla est mattis ut. Nunc lobortis ullamcorper elit, quis tincidunt lacus. Quisque ac risus mauris.

Here is an ordered list of outputs:
1. This is item 1.
2. This is item 2.
3. This is item 3.
4. This is item 4.
5. This is item 5.
6. This is item 6.

Here is an unordered list:
- This is item 1.
- This is item 2.
- This is item 3.
- This is item 4.
- This is item 5.
- This is item 6.

### Tersary Title

Fusce vehicula justo vitae lectus porttitor porttitor. Aenean id augue eu lacus laoreet mollis. Pellentesque porta nunc at libero aliquet pharetra. Mauris nec ipsum a sapien faucibus commodo. Nulla vel vehicula diam. Proin elementum dolor at sapien dapibus, id molestie enim rutrum. Nunc semper ligula lorem, in luctus justo efficitur sit amet. Morbi purus neque, volutpat quis efficitur sed, pretium et turpis. Aenean vehicula viverra viverra. Phasellus nec interdum risus. In efficitur, mi a aliquet faucibus, est arcu convallis turpis, id tincidunt mauris purus at arcu. Pellentesque nisi velit, dictum et turpis sed, blandit fermentum ante. Pellentesque vestibulum sapien vitae diam lacinia, nec dignissim dolor gravida. Maecenas fermentum sed mauris sit amet placerat. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos.

### Tersary Title

Pellentesque vitae tellus elit. Nunc faucibus tincidunt ante sit amet porta. Vestibulum imperdiet mauris pretium arcu tincidunt, eu rhoncus lectus condimentum. Nulla bibendum porta augue, vel facilisis ipsum sodales id. Phasellus interdum commodo nisl sed elementum. Nullam in rutrum ex, at euismod sapien. Aliquam ullamcorper dolor a enim sagittis, quis varius ipsum sodales. Nullam vel tincidunt lectus, eget placerat felis. Sed ac quam mi. Vestibulum et tortor id ante viverra feugiat. In mollis venenatis massa eget vehicula. Nam pretium risus elit, eu porttitor sapien tempor ut.
