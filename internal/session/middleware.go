package session

import "net/http"

func (n *Session) SessionMiddleware(next http.Handler) http.Handler {
	return n.session.LoadAndSave(next)
}
