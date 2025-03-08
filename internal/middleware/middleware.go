package middleware

import (
	"net/http"
	"time"

	"github.com/PsionicAlch/course-platform/internal/authentication"
	"github.com/PsionicAlch/course-platform/internal/render"
	"github.com/PsionicAlch/course-platform/web/html"
	"github.com/go-chi/httprate"
	"github.com/justinas/nosurf"
)

// CSRFProtection is middleware to set up nosurf CSRF protection.
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

// RateLimiter is middleware to set up httprate rate limiting.
func RateLimiter(requestLimit int, windowLength time.Duration, renderer render.Renderer) func(next http.Handler) http.Handler {
	return httprate.Limit(
		requestLimit,
		windowLength,
		httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
			user := authentication.GetUserFromRequest(r)

			if user != nil {
				return user.ID, nil
			}

			return httprate.KeyByRealIP(r)
		}),
		httprate.WithLimitHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := authentication.GetUserFromRequest(r)
			pageData := &html.Errors429Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}

			renderer.RenderHTML(w, nil, "errors-429", pageData)
		})),
		httprate.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
			user := authentication.GetUserFromRequest(r)
			pageData := &html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}

			renderer.RenderHTML(w, nil, "errors-500", pageData)
		}),
	)
}
