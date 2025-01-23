package session

import "net/http"

// TODO: Write documentation

func (n *Session) SessionMiddleware(next http.Handler) http.Handler {
	return n.session.LoadAndSave(next)
}
