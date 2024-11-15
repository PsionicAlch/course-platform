package authentication

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
)

func (auth *Authentication) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetUserFromAuthCookie(r.Cookies())
		if err != nil {
			auth.ErrorLog.Printf("Failed to get user from cookies: %s\n", err)
			user = nil
		}

		ctx := NewContextWithUserModel(user, r.Context())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (auth *Authentication) AllowAuthenticated(redirectURL string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromRequest(r)

			if user == nil {
				utils.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (auth *Authentication) AllowUnauthenticated(redirectURL string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromRequest(r)
			if user != nil {
				utils.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (auth *Authentication) AllowAdmin(redirectURL string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromRequest(r)
			if user == nil {
				utils.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				return
			}

			if !user.IsAdmin {
				utils.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
