package middleware

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func CSRFProtection(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	return csrfHandler
}
