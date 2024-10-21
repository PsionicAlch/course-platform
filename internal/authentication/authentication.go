package authentication

import (
	"fmt"

	"github.com/gorilla/securecookie"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/forms"
)

type Authentication struct {
	db           database.Database
	secureCookie *securecookie.SecureCookie
}

func CreateAuthentication(db database.Database) *Authentication {
	s := securecookie.New()
	return &Authentication{
		db: db,
	}
}

func (auth *Authentication) SignUserIn(form *forms.SignUpForm) error {
	// Hash user's password.
	password, err := HashPassword(form.GetPassword())
	if err != nil {
		return CreateFailedToSignUserIn(fmt.Sprintf("failed to hash user's password: %s", err))
	}

	// Add user to database.
	userId, err := auth.db.AddUser(form.GetEmail(), password)
	if err != nil {
		return CreateFailedToSignUserIn(err.Error())
	}

	// Generate authentication token.
	token, err := auth.db.CreateAuthenticationToken(userId)
	if err != nil {
		return CreateFailedToSignUserIn(err.Error())
	}

	// Save authentication token to secure cookie.

	return nil
}
