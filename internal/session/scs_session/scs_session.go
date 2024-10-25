package scssession

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/config"
	"github.com/PsionicAlch/psionicalch-home/website/forms"
	"github.com/alexedwards/scs/v2"
)

type SCSSession struct {
	SessionManager *scs.SessionManager
}

func NewSession() *SCSSession {
	// Set up *scs.SessionManager.
	sessionManager := scs.New()
	sessionManager.Lifetime = 3 * time.Hour
	sessionManager.IdleTimeout = 20 * time.Minute
	sessionManager.Cookie.Name = config.GetWithoutError[string]("SESSION_COOKIE_NAME")
	sessionManager.Cookie.Domain = config.GetWithoutError[string]("DOMAIN_NAME")
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Path = "/"
	sessionManager.Cookie.Persist = false
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = config.InProduction()

	// Register the required data types for retrieval.
	gob.Register(&forms.SignUpForm{})
	gob.Register(&forms.LoginForm{})

	return &SCSSession{
		SessionManager: sessionManager,
	}
}
