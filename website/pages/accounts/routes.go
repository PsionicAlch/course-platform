package accounts

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Use(handlers.auth.AllowUnauthenticated("/profile"))

	router.Get("/login", handlers.LoginGet)
	router.Post("/login", handlers.LoginPost)

	router.Get("/signup", handlers.SignupGet)
	router.Post("/signup", handlers.SignupPost)

	router.Get("/forgot", handlers.ForgotGet)

	router.Post("/validate/signup-form", handlers.ValidateSignupForm)

	return router
}
