package middleware

import (
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/web/html"
	"github.com/go-chi/httprate"
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
