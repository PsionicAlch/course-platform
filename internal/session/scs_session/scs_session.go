package scssession

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/config"
	"github.com/PsionicAlch/psionicalch-home/internal/forms"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/alexedwards/scs/v2"
)

type SCSSession struct {
	utils.Loggers
	SessionManager *scs.SessionManager
}

func NewSession() *SCSSession {
	// Set up loggers.
	sessionLoggers := utils.CreateLoggers("SESSION")

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

	return &SCSSession{
		Loggers:        sessionLoggers,
		SessionManager: sessionManager,
	}
}

func (s *SCSSession) StoreSignUpFormData(ctx context.Context, signUpFormData *forms.SignUpForm) {
	s.SessionManager.Put(ctx, "signup-form-data", signUpFormData)
}

func (s *SCSSession) RetrieveSignUpFormData(ctx context.Context) *forms.SignUpForm {
	signUpFormData := forms.NewSignupForm()

	if s.SessionManager.Exists(ctx, "signup-form-data") {
		data, ok := s.SessionManager.Pop(ctx, "signup-form-data").(*forms.SignUpForm)
		if !ok {
			s.ErrorLog.Println("Failed to retrieve SignUp Form Data from session!")
			return signUpFormData
		}

		signUpFormData = data
	}

	return signUpFormData
}
