package profile

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	// router.Use(handlers.auth.AllowAuthenticated("/accounts/login"))

	router.Get("/", handlers.ProfileGet)

	router.Get("/affiliate-history", handlers.AffiliateHistoryGet)

	router.Get("/courses", handlers.CoursesGet)
	router.Get("/courses/{slug}", handlers.CourseGet)
	router.Get("/courses/{course_slug}/{chapter_slug}", handlers.CourseChapterGet)

	router.Get("/tutorials/bookmarks", handlers.TutorialsBookmarksGet)
	router.Get("/tutorials/liked", handlers.TutorialsLikedGet)

	return router
}
