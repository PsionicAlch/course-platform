package database

import (
	"math/rand"
	"time"

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
