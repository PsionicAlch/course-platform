package database

import (
	cryptoRand "crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	mathRand "math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

// GenerateID generates a new ULID to be used as an ID.
func GenerateID() (string, error) {
	now := time.Now()
	entropy := mathRand.New(mathRand.NewSource(now.UnixNano()))
	ms := ulid.Timestamp(now)
	id, err := ulid.New(ms, entropy)

	return id.String(), err
}

// GenerateToken creates a new token for use in authentication, emails, and payments.
func GenerateToken() (string, error) {
	tokenBytes, err := RandomBytes(32)
	if err != nil {
		return "", err
	}

	return BytesToURLString(tokenBytes), nil
}

// RandomBytes generates a slice of random bytes equal to the given length.
func RandomBytes(length uint) ([]byte, error) {
	b := make([]byte, length)
	_, err := cryptoRand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// BytesToURLString converts a slice of bytes to a URL safe string.
func BytesToURLString(src []byte) string {
	return base64.RawURLEncoding.EncodeToString(src)
}

// NewNullString creates a new sql.NullString based off a given string. If the string is empty the sql.NullString
// generated won't be valid.
func NewNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{String: s, Valid: true}
}

// NameSurnameToSlug takes a user's name and surname then returns
// a URL safe slug that can be used to find the user in the frontend.
// I didn't want to unnecessarily expose the user's ID.
//
// Keep in mind that this function was written by ChatGPT 4o.
func NameSurnameToSlug(name, surname string) string {
	// Combine name and surname with a hyphen
	fullName := fmt.Sprintf("%s-%s", name, surname)

	// Convert to lowercase
	fullName = strings.ToLower(fullName)

	// Remove non-alphanumeric characters except for hyphens
	re := regexp.MustCompile(`[^a-z0-9-]+`)
	safeSlug := re.ReplaceAllString(fullName, "")

	// Trim leading and trailing hyphens
	safeSlug = strings.Trim(safeSlug, "-")

	return safeSlug
}
