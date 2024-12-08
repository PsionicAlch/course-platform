package authentication

import (
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

const (
	AuthenticationToken = "authentication"
	EmailToken          = "email"
)

// TODO: Move token logic over to using database.GenerateToken

func NewToken() (string, error) {
	tokenBytes, err := RandomBytes(32)
	if err != nil {
		return "", err
	}

	return BytesToURLString(tokenBytes), nil
}

func ValidateToken(token *models.TokenModel, tokenType string) bool {
	if token == nil {
		return false
	}

	if token.Token == "" || token.TokenType != tokenType {
		return false
	}

	if time.Now().After(token.ValidUntil) {
		return false
	}

	return true
}
