package scssession

import (
	"context"

	"github.com/PsionicAlch/psionicalch-home/website/forms"
)

// StoreSignUpFormData stores an instance of the signup form in the session.
func (s *SCSSession) StoreLoginFormData(ctx context.Context, signUpFormData *forms.LoginForm) {
	s.SessionManager.Put(ctx, "login-form-data", signUpFormData)
}

// RetrieveSignUpFormData retrieves an instance of the signup form from the session. If there is no
// sign up form currently stored in the session it will return a freshly created instance.
func (s *SCSSession) RetrieveLoginFormData(ctx context.Context) *forms.LoginForm {
	signUpFormData := forms.NewLoginForm()

	if s.SessionManager.Exists(ctx, "signup-form-data") {
		data, ok := s.SessionManager.Pop(ctx, "login-form-data").(*forms.LoginForm)
		if !ok {
			return signUpFormData
		}

		signUpFormData = data
	}

	return signUpFormData
}
