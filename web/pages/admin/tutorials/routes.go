package tutorials

import (
	"net/http"

	"github.com/PsionicAlch/course-platform/web/pages"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Get("/", handlers.TutorialsGet)
	router.Get("/htmx", handlers.TutorialsPaginationGet)

	router.Get("/htmx/change-published/{tutorial-id}", handlers.PublishedEditGet)
	router.Post("/htmx/change-published/{tutorial-id}", handlers.PublishedEditPost)

	router.Get("/htmx/change-author/{tutorial-id}", handlers.AuthorEditGet)
	router.Post("/htmx/change-author/{tutorial-id}", handlers.AuthorEditPost)

	return router
}
