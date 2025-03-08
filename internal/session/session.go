package session

import (
	"context"
	"net/http"
	"time"

	"github.com/PsionicAlch/course-platform/internal/utils"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

// TODO: Write documentation

type Session struct {
	utils.Loggers
	session *scs.SessionManager
}

func SetupSession(cookieName, cookieDomain string) *Session {
	loggers := utils.CreateLoggers("SESSION")

	session := scs.New()
	session.Lifetime = 1 * time.Hour
	session.IdleTimeout = 20 * time.Minute
	session.Cookie.Name = cookieName
	session.Cookie.Domain = cookieDomain
	session.Cookie.HttpOnly = true
	session.Cookie.Path = "/"
	session.Cookie.Persist = false
	session.Cookie.SameSite = http.SameSiteStrictMode
	session.Cookie.Secure = true
	session.Store = memstore.New()

	return &Session{
		Loggers: loggers,
		session: session,
	}
}

func (s *Session) Reset(ctx context.Context) {
	if err := s.session.Destroy(ctx); err != nil {
		s.ErrorLog.Printf("Failed to destroy the current session: %s\n", err)
	}
}
