package gatekeeper

import (
	"errors"
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

func NewToken(token, tokenType, userId, ipAddr string, validUntil time.Time) (*Token, error) {
	// TODO: Custom errors.
	if token == "" {
		return nil, errors.New("token cannot be empty")
	}

	if tokenType == "" {
		return nil, errors.New("token type cannot be empty")
	}

	if !inSlice(tokenType, []string{authenticationTokenType, emailTokenType}) {
		return nil, errors.New("invalid token type")
	}

	if userId == "" {
		return nil, errors.New("user id cannot be empty")
	}

	newToken := &Token{
		Token:      token,
		TokenType:  tokenType,
		ValidUntil: validUntil,
		UserID:     userId,
		IPAddress:  ipAddr,
	}

	return newToken, nil
}
