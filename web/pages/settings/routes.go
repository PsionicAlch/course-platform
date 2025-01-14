package settings

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/web/pages"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(handlerContext *pages.HandlerContext) http.Handler {
	handlers := SetupHandlers(handlerContext)

	router := chi.NewRouter()

	router.Use(handlerContext.Authentication.SetUserWithEmail(handlers.Emailer))
	router.Use(handlerContext.Session.SessionMiddleware)
	router.Use(handlerContext.Authentication.AllowAuthenticated("/accounts/login"))

	router.Get("/", handlers.SettingsGet)

	router.Get("/whitelist/{ip-address}", handlers.WhitelistIPAddressPost)

	router.Post("/change-first-name", handlers.ChangeFirstNamePost)

	router.Post("/change-last-name", handlers.ChangeLastNamePost)

	router.Post("/change-email", handlers.ChangeEmailPost)

	router.Post("/change-password", handlers.ChangePasswordPost)

	router.Delete("/delete-ip-address/{ip-address-id}", handlers.IPAddressDelete)

	router.Post("/request-refund/{course-id}", handlers.RequestRefundPost)

	router.Delete("/delete-account", handlers.AccountDelete)

	router.Post("/validate/change-password", handlers.ValidateChangePassword)

	return router
}
