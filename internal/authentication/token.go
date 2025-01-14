package authentication

import (
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

const (
	AuthenticationToken = "authentication"
	EmailToken          = "email"
)

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
