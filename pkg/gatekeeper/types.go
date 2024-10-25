package gatekeeper

import (
	"time"
)

const (
	authenticationTokenType = "authentication"
	emailTokenType          = "email"
)

type GatekeeperPassword struct {
	ArgonVersion int
	Hash         []byte
	Salt         []byte
	Iterations   uint8
	Memory       uint32
	Threads      uint8
	KeyLength    uint8
}

type Token struct {
	Token      string
	TokenType  string
	ValidUntil time.Time
	UserID     string
	IPAddress  string
}

func NewToken(token, tokenType, userId, ipAddr string, validUntil time.Time) *Token {
	return &Token{
		Token:      token,
		TokenType:  tokenType,
		ValidUntil: validUntil,
		UserID:     userId,
		IPAddress:  ipAddr,
	}
}

type UserInformation struct {
	ID       string
	Email    string
	Password string
}

func NewUser(id, email, password string) *UserInformation {
	return &UserInformation{
		ID:       id,
		Email:    email,
		Password: password,
	}
}
