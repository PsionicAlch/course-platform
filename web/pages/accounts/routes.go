package accounts

import (
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/middleware"
	"github.com/PsionicAlch/psionicalch-home/web/pages"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Use(handlerContext.Authentication.SetUser)
	router.Use(handlerContext.Session.SessionMiddleware)

	router.With(handlerContext.Authentication.AllowAuthenticated("/accounts/login")).Delete("/logout", handlers.LogoutDelete)

	router.With(handlerContext.Authentication.AllowUnauthenticated("/profile")).Get("/login", handlers.LoginGet)
	router.With(
		handlerContext.Authentication.AllowUnauthenticated("/profile"),
		middleware.RateLimiter(5, time.Minute, handlerContext.Renderers.Page),
	).Post("/login", handlers.LoginPost)

	router.With(handlerContext.Authentication.AllowUnauthenticated("/profile")).Get("/signup", handlers.SignupGet)
	router.With(
		handlerContext.Authentication.AllowUnauthenticated("/profile"),
		middleware.RateLimiter(5, time.Minute, handlerContext.Renderers.Page),
	).Post("/signup", handlers.SignupPost)

	router.With(handlerContext.Authentication.AllowUnauthenticated("/settings#change-password")).Get("/reset-password", handlers.ForgotPasswordGet)
	router.With(
		handlerContext.Authentication.AllowUnauthenticated("/settings#change-password"),
		middleware.RateLimiter(5, time.Minute, handlerContext.Renderers.Page),
	).Post("/reset-password", handlers.ForgotPasswordPost)

	router.With(handlerContext.Authentication.AllowUnauthenticated("/settings#change-password")).Get("/reset-password/{email_token}", handlers.ResetPasswordGet)
	router.With(
		handlerContext.Authentication.AllowUnauthenticated("/settings#change-password"),
		middleware.RateLimiter(5, time.Minute, handlerContext.Renderers.Page),
	).Post("/reset-password/{email_token}", handlers.ResetPasswordPost)

	router.With(handlerContext.Authentication.AllowUnauthenticated("/")).Post("/validate/signup", handlers.ValidateSignupPost)
	router.With(handlerContext.Authentication.AllowUnauthenticated("/")).Post("/validate/reset-password/{email_token}", handlers.ValidateResetPasswordPost)

	return router
}
