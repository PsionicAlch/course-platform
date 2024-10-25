package gatekeeper

import (
	"crypto/subtle"
	"errors"

	"golang.org/x/crypto/argon2"
)

type GatekeeperPasswordHashParameters struct {
	saltLength uint8
	iterations uint8
	memory     uint32
	threads    uint8
	keyLength  uint8
}

// CreatePasswordHashParameters creates a new instance of the GatekeeperPasswordHashParameters.
func CreatePasswordHashParameters(saltLength uint8, iterations uint8, memory uint32, threads uint8, keyLength uint8) *GatekeeperPasswordHashParameters {
	return &GatekeeperPasswordHashParameters{
		saltLength: saltLength,
		iterations: iterations,
		memory:     memory,
		threads:    threads,
		keyLength:  keyLength,
	}
}

func (params *GatekeeperPasswordHashParameters) HashPassword(password string) (string, error) {
	salt, err := newSalt(int(params.saltLength))
	if err != nil {
		// TODO: Create dedicated error.
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, uint32(params.iterations), params.memory, params.threads, uint32(params.keyLength))

	passwordStruct := GatekeeperPassword{
		ArgonVersion: argon2.Version,
		Hash:         hash,
		Salt:         salt,
		Iterations:   params.iterations,
		Memory:       params.memory,
		Threads:      params.threads,
		KeyLength:    params.keyLength,
	}

	passwordBytes, err := passwordToBytes(&passwordStruct)
	if err != nil {
		// TODO: Create dedicated error.
		return "", err
	}

	return bytesToString(passwordBytes), nil
}

func ComparePasswordAndHash(password, passwordHash string) (bool, error) {
	passwordBytes, err := stringToBytes(passwordHash)
	if err != nil {
		// TODO: Create dedicated error.
		return false, err
	}

	passwordStruct, err := passwordFromBytes(passwordBytes)
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
