package database

import "time"

type UserModel struct {
	ID         string
	Email      string
	Password   string
	Created_At time.Time
	Update_At  time.Time
}

type TokenModel struct {
	ID         string
	Token      string
	TokenType  string
	ValidUntil time.Time
	CreatedAt  time.Time
	UserID     string
}
