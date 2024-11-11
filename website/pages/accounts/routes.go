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

	router.Get("/reset-password", handlers.ForgotGet)
	router.Get("/reset-password/{email_token}", handlers.ResetPasswordGet)

	router.Post("/validate/signup-form", handlers.ValidateSignupForm)

	return router
}
