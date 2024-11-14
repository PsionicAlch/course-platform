package authentication

import "net/http"

func (auth *Authentication) SetUserMiddleware(next http.Handler) http.Handler {
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
