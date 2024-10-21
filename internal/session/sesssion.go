package session

import (
	"context"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/forms"
)

type Session interface {
	StoreSignUpFormData(ctx context.Context, signUpFormData *forms.SignUpForm)
	RetrieveSignUpFormData(ctx context.Context) *forms.SignUpForm

	LoadSession(next http.Handler) http.Handler
}
