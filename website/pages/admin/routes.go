package admin

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.AdminGet)
	router.Get("/admins", handlers.AdminsGet)
	router.Get("/authors", handlers.AuthorsGet)
	router.Get("/comments", handlers.CommentsGet)
	router.Get("/courses", handlers.CoursesGet)
	router.Get("/discounts", handlers.DiscountsGet)
	router.Get("/purchases", handlers.PurchasesGet)
	router.Get("/refunds", handlers.RefundsGet)
	router.Get("/tutorials", handlers.TutorialsGet)
	router.Get("/users", handlers.UsersGet)

	return router
}
