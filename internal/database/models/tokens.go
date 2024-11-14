package models

import "time"

type TokenModel struct {
	ID         string
	Token      string
	TokenType  string
	ValidUntil time.Time
	UserID     string
	IPAddr     string
	CreatedAt  time.Time
}
