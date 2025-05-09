package users

import (
	"net/http"

	"github.com/PsionicAlch/course-platform/web/pages"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Get("/", handlers.UsersGet)
	router.Get("/htmx", handlers.UsersPaginationGet)

	router.Get("/htmx/change-author/{user-id}", handlers.AuthorEditGet)
	router.Post("/htmx/change-author/{user-id}", handlers.AuthorEditPost)

	router.Get("/htmx/change-admin/{user-id}", handlers.AdminEditGet)
	router.Post("/htmx/change-admin/{user-id}", handlers.AdminEditPost)

	return router
}
