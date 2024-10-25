package session

import (
	"context"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/website/forms"
)

type Session interface {
	// SignUpForm based functions.
	StoreSignUpFormData(ctx context.Context, signUpFormData *forms.SignUpForm)
	RetrieveSignUpFormData(ctx context.Context) *forms.SignUpForm

	// LoginForm based functions.
	StoreLoginFormData(ctx context.Context, loginFormData *forms.LoginForm)
	RetrieveLoginFormData(ctx context.Context) *forms.LoginForm

	// Middleware functions.
	LoadSession(next http.Handler) http.Handler
}
