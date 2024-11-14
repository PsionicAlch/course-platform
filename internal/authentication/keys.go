package authentication

import (
	"fmt"
	"strings"
)

type SecureCookieKeys struct {
	HashKey  []byte
	BlockKey []byte
}

func CreateSecureCookieKeys(key string) (*SecureCookieKeys, error) {
	// This makes things easier to work with if we don't yet have a
	// previous key.
	if key == "" {
		return new(SecureCookieKeys), nil
	}

	hashKey, blockKey, found := strings.Cut(key, "$")
	if !found {
		return nil, ErrInvalidSecureCookieKey
	}

	hashKeyBytes, err := StringToBytes(hashKey)
	if err != nil {
		return nil, ErrInvalidSecureCookieKey
	}

	blockKeyBytes, err := StringToBytes(blockKey)
	if err != nil {
		return nil, ErrInvalidSecureCookieKey
	}

	secureCookieKeys := &SecureCookieKeys{
		HashKey:  hashKeyBytes,
		BlockKey: blockKeyBytes,
	}

	return secureCookieKeys, nil
}

func GenerateKeyString() (string, error) {
	hashKey, err := GenerateKey(64)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash key: %w", err)
	}

	blockKey, err := GenerateKey(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate block key: %w", err)
	}

	return fmt.Sprintf("%s$%s", hashKey, blockKey), nil
}

func GenerateKey(length uint) (string, error) {
	byteSlice, err := RandomBytes(length)
	if err != nil {
		return "", err
	}

	return BytesToString(byteSlice), nil
}
