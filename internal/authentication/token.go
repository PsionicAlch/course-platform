package authentication

import (
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

const (
	AuthenticationToken = "authentication"
	EmailToken          = "email"
)

func NewToken() (string, error) {
	tokenBytes, err := RandomBytes(32)
	if err != nil {
		return "", err
	}

	return BytesToString(tokenBytes), nil
}

func ValidateAuthenticationToken(token *models.TokenModel) bool {
	if token == nil {
		return false
	}

	if token.Token == "" || token.TokenType != AuthenticationToken {
		return false
	}

	if time.Now().After(token.ValidUntil) {
		return false
	}

	return true
}
