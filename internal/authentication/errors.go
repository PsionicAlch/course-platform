package authentication

import "errors"

var (
	ErrUserExists             = errors.New("user already exists")
	ErrEmptySecureCookieKey   = errors.New("secure cookies key is empty")
	ErrInvalidSecureCookieKey = errors.New("invalid secure cookies key")
)
