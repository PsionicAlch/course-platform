package accounts

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.With(handlers.Auth.AllowAuthenticated("/accounts/login")).Delete("/logout", handlers.LogoutDelete)

	router.With(handlers.Auth.AllowUnauthenticated("/profile")).Get("/login", handlers.LoginGet)
	router.With(handlers.Auth.AllowUnauthenticated("/profile")).Post("/login", handlers.LoginPost)

	router.With(handlers.Auth.AllowUnauthenticated("/profile")).Get("/signup", handlers.SignupGet)
	router.With(handlers.Auth.AllowUnauthenticated("/profile")).Post("/signup", handlers.SignupPost)

	router.With(handlers.Auth.AllowUnauthenticated("/profile")).Get("/reset-password", handlers.ForgotGet)

	router.With(handlers.Auth.AllowUnauthenticated("/profile")).Get("/reset-password/{email_token}", handlers.ResetPasswordGet)

	router.With(handlers.Auth.AllowUnauthenticated("/profile")).Post("/validate/signup", handlers.ValidateSignupPost)

	return router
}
