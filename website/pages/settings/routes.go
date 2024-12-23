package settings

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Get("/", handlers.SettingsGet)

	router.Post("/change-first-name", handlers.ChangeFirstNamePost)

	router.Post("/change-last-name", handlers.ChangeLastNamePost)

	router.Post("/change-email", handlers.ChangeEmailPost)

	router.Post("/change-password", handlers.ChangePasswordPost)

	router.Post("/validate/change-password", handlers.ValidateChangePassword)

	return router
}
