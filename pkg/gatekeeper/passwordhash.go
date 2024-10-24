package gatekeeper

import (
	"crypto/subtle"
	"errors"

	"golang.org/x/crypto/argon2"
)

type GatekeeperPasswordHashParameters struct {
	SaltLength uint8
	Iterations uint8
	Memory     uint32
	Threads    uint8
	KeyLength  uint8
}

// CreatePasswordHashParameters creates a new instance of the GatekeeperPasswordHashParameters.
func CreatePasswordHashParameters(saltLength uint8, iterations uint8, memory uint32, threads uint8, keyLength uint8) *GatekeeperPasswordHashParameters {
	return &GatekeeperPasswordHashParameters{
		SaltLength: saltLength,
		Iterations: iterations,
		Memory:     memory,
		Threads:    threads,
		KeyLength:  keyLength,
	}
}

func (params *GatekeeperPasswordHashParameters) HashPassword(password string) (string, error) {
	salt, err := NewSalt(int(params.SaltLength))
	if err != nil {
		// TODO: Create dedicated error.
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, uint32(params.Iterations), params.Memory, params.Threads, uint32(params.KeyLength))

	passwordStruct := GatekeeperPassword{
		ArgonVersion: argon2.Version,
		Hash:         hash,
		Salt:         salt,
		Iterations:   params.Iterations,
		Memory:       params.Memory,
		Threads:      params.Threads,
		KeyLength:    params.KeyLength,
	}

	passwordBytes, err := PasswordToBytes(&passwordStruct)
	if err != nil {
		// TODO: Create dedicated error.
		return "", err
	}

	return BytesToString(passwordBytes), nil
}

func ComparePasswordAndHash(password, passwordHash string) (bool, error) {
	passwordBytes, err := StringToBytes(passwordHash)
	if err != nil {
		// TODO: Create dedicated error.
		return false, err
	}

	passwordStruct, err := PasswordFromBytes(passwordBytes)
	if err != nil {
		// TODO: Create dedicated error.
		return false, err
	}

	if passwordStruct.ArgonVersion != argon2.Version {
		// TODO: Create dedicated error.
		return false, errors.New("mismatched argon versions between package and hash")
	}

	newPasswordHash := argon2.IDKey([]byte(password), passwordStruct.Salt, uint32(passwordStruct.Iterations), passwordStruct.Memory, passwordStruct.Threads, uint32(passwordStruct.KeyLength))

	if subtle.ConstantTimeCompare(passwordStruct.Hash, newPasswordHash) == 1 {
		return true, nil
	}

	return false, nil
}
