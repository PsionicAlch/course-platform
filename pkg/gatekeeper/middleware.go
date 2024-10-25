package gatekeeper

import "net/http"

func (gatekeeper *Gatekeeper) AllowAuthenticated(loginUrl string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if authenticated, err := gatekeeper.ValidateAuthenticationToken(r.Cookies()); err != nil || !authenticated {
				redirect(w, r, loginUrl)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (gatekeeper *Gatekeeper) AllowUnauthenticated(loggedInUrl string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if authenticated, err := gatekeeper.ValidateAuthenticationToken(r.Cookies()); err == nil && authenticated {
				redirect(w, r, loggedInUrl)
			}

			next.ServeHTTP(w, r)
		})
	}
}
