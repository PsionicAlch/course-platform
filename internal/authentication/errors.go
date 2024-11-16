package authentication

import "errors"

var (
	ErrUserExists             = errors.New("user already exists")
	ErrEmptySecureCookieKey   = errors.New("secure cookies key is empty")
	ErrInvalidSecureCookieKey = errors.New("invalid secure cookies key")
	ErrInvalidCredentials     = errors.New("invalid user credentials")
	ErrMismatchedArgonVersion = errors.New("mismatched argon 2 version")
	ErrUnregisteredEmail      = errors.New("no account is registered with this email")
)
