package scssession

import (
	"context"

	"github.com/PsionicAlch/psionicalch-home/internal/forms"
)

// StoreSignUpFormData stores an instance of the signup form in the session.
func (s *SCSSession) StoreSignUpFormData(ctx context.Context, signUpFormData *forms.SignUpForm) {
	s.SessionManager.Put(ctx, "signup-form-data", signUpFormData)
}

// RetrieveSignUpFormData retrieves an instance of the signup form from the session. If there is no
// sign up form currently stored in the session it will return a freshly created instance.
func (s *SCSSession) RetrieveSignUpFormData(ctx context.Context) *forms.SignUpForm {
	signUpFormData := forms.NewSignupForm()

	if s.SessionManager.Exists(ctx, "signup-form-data") {
		data, ok := s.SessionManager.Pop(ctx, "signup-form-data").(*forms.SignUpForm)
		if !ok {
			return signUpFormData
		}

		signUpFormData = data
	}

	return signUpFormData
}
