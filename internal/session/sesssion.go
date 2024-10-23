package session

import (
	"context"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/forms"
)

type Session interface {
	// SignUpForm based functions.
	StoreSignUpFormData(ctx context.Context, signUpFormData *forms.SignUpForm)
	RetrieveSignUpFormData(ctx context.Context) *forms.SignUpForm

	// Middleware functions.
	LoadSession(next http.Handler) http.Handler
}
