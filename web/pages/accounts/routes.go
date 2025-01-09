package accounts

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlers *Handlers) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.CSRFProtection)

	router.With(handlers.Auth.AllowAuthenticated("/accounts/login")).Delete("/logout", handlers.LogoutDelete)

	router.With(handlers.Auth.AllowUnauthenticated("/profile")).Get("/login", handlers.LoginGet)
	router.With(handlers.Auth.AllowUnauthenticated("/profile")).Post("/login", handlers.LoginPost)

	router.With(handlers.Auth.AllowUnauthenticated("/profile")).Get("/signup", handlers.SignupGet)
	router.With(handlers.Auth.AllowUnauthenticated("/profile")).Post("/signup", handlers.SignupPost)

	router.With(handlers.Auth.AllowUnauthenticated("/settings#change-password")).Get("/reset-password", handlers.ForgotPasswordGet)
	router.With(handlers.Auth.AllowUnauthenticated("/settings#change-password")).Post("/reset-password", handlers.ForgotPasswordPost)

	router.With(handlers.Auth.AllowUnauthenticated("/settings#change-password")).Get("/reset-password/{email_token}", handlers.ResetPasswordGet)
	router.With(handlers.Auth.AllowUnauthenticated("/settings#change-password")).Post("/reset-password/{email_token}", handlers.ResetPasswordPost)

	router.With(handlers.Auth.AllowUnauthenticated("/")).Post("/validate/signup", handlers.ValidateSignupPost)
	router.With(handlers.Auth.AllowUnauthenticated("/")).Post("/validate/reset-password/{email_token}", handlers.ValidateResetPasswordPost)

	return router
}
