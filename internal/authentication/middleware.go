package authentication

// func (auth *Authentication) AllowAuthenticated(redirectURL string) func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			user, err := auth.GetUserFromCookie(r.Cookies())
// 			if err != nil {
// 				auth.Loggers.WarningLog.Println("Failed to get user from authentication cookie: ", err)

// 				utils.Redirect(w, r, "/accounts/login", http.StatusTemporaryRedirect)
// 				return
// 			}

// 			if user == nil {
// 				utils.Redirect(w, r, "/accounts/login", http.StatusTemporaryRedirect)
// 				return
// 			}

// 			auth.Loggers.InfoLog.Printf("AllowAuthenticated middleware was called!\nHere is the cookies: %#v\n", r.Cookies())
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }
