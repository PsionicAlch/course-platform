package database

import (
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/oklog/ulid/v2"
)

// GenerateID generates a new ULID to be used as an ID.
func GenerateID() (string, error) {
	now := time.Now()
	entropy := rand.New(rand.NewSource(now.UnixNano()))
	ms := ulid.Timestamp(now)
	id, err := ulid.New(ms, entropy)

	return id.String(), err
}

// GenerateToken generates a new token.
func GenerateToken() (string, error) {
	tokenBytes, err := utils.RandomByteSlice(32)
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(tokenBytes), nil
}
