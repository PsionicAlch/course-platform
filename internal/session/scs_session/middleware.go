package scssession

import "net/http"

func (s *SCSSession) LoadSession(next http.Handler) http.Handler {
	return s.SessionManager.LoadAndSave(next)
}
