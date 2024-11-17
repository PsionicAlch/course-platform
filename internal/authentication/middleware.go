package authentication

import (
	"net"
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
)

type AuthenticationEmail interface {
	SendSuspiciousActivityEmail(email, firstName, ipAddr string, dateTime time.Time)
}

func (auth *Authentication) SetUserWithEmail(email AuthenticationEmail) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, token, err := auth.GetUserFromAuthCookie(r.Cookies())
			if err != nil {
				auth.ErrorLog.Printf("Failed to get user from cookies: %s\n", err)
				user = nil
			}

			// In case this account was accessed from a different IP address send the user
			// so that they can take any actions required to secure their account.
			ipAddr, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				auth.ErrorLog.Printf("Failed to get ip address from r.RemoteAddr: %s\n", err)
			} else {
				if user != nil && ipAddr != token.IPAddr {
					go email.SendSuspiciousActivityEmail(user.Email, user.Name, ipAddr, time.Now())
				}
			}

			ctx := NewContextWithUserModel(user, r.Context())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (auth *Authentication) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _, err := auth.GetUserFromAuthCookie(r.Cookies())
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
