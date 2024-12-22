package authentication

import (
	"net"
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
)

type AuthenticationEmail interface {
	SendSuspiciousActivityEmail(email, firstName, ipAddr string, dateTime time.Time)
}

func (auth *Authentication) SetUserWithEmail(email AuthenticationEmail) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := auth.GetUserFromAuthCookie(r.Cookies())
			if err != nil {
				auth.ErrorLog.Printf("Failed to get user from cookies: %s\n", err)
				user = nil
			}

			// In case this account was accessed from a IP address that's not in their
			// whitelist we want to notify them of the suspicious behavior.
			if user != nil {
				ipAddr, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					auth.ErrorLog.Printf("Failed to get ip address from r.RemoteAddr: %s\n", err)
				}

				userIpAddresses, err := auth.Database.GetUserIpAddresses(user.ID)
				if err != nil {
					auth.ErrorLog.Printf("Failed to get user's (\"%s\") whitelisted IP addresses: %s\n", user.Email, err)
				}

				_, whitelistedIP := utils.InSliceFunc(ipAddr, userIpAddresses, func(ip string, addr *models.WhitelistedIPModel) bool {
					return ip == addr.IPAddress
				})

				if userIpAddresses != nil && !whitelistedIP {
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
				if redirectURL != "" {
					auth.Session.SetRedirectURL(r.Context(), r.URL.Path)
					utils.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				}

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
				if redirectURL != "" {
					auth.Session.SetRedirectURL(r.Context(), r.URL.Path)
					utils.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				}

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
				if redirectURL != "" {
					auth.Session.SetRedirectURL(r.Context(), r.URL.Path)
					utils.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				}

				return
			}

			if !user.IsAdmin {
				if redirectURL != "" {
					auth.Session.SetRedirectURL(r.Context(), r.URL.Path)
					utils.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				}

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
