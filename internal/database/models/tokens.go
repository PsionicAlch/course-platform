package models

import "time"

// TokenModel is a struct representation of the tokens table.
type TokenModel struct {
	ID         string
	Token      string
	TokenType  string
	ValidUntil time.Time
	UserID     string
	CreatedAt  time.Time
}
