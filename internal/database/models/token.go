package models

import "time"

type TokenModel struct {
	ID         string
	Token      string
	TokenType  string
	ValidUntil time.Time
	CreatedAt  time.Time
	UserID     string
}
