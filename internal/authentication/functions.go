package authentication

import (
	"bytes"
	"crypto/subtle"
	"encoding/gob"
	"errors"
	"runtime"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"golang.org/x/crypto/argon2"
)

type HashParameters struct {
	SaltLength uint8
	Iterations uint8
	Memory     uint32
	Threads    uint8
	KeyLength  uint8
}

type Password struct {
	ArgonVersion int
	Hash         []byte
	Salt         []byte
	Iterations   uint8
	Memory       uint32
	Threads      uint8
	KeyLength    uint8
}

// HashPassword takes a password string and some optional parameters and returns a hashed password string.
func HashPassword(password string, params ...*HashParameters) (string, error) {
	// Create hash parameters.
	hashParams := new(HashParameters)
	if len(params) > 0 {
		hashParams = params[0]
	} else {
		hashParams.SaltLength = 16
		hashParams.Iterations = uint8(runtime.NumCPU())
		hashParams.Memory = 64 * 1024
		hashParams.Threads = uint8(runtime.NumCPU())
		hashParams.KeyLength = 32
	}

	// Create a salt byte slice.
	salt, err := utils.RandomByteSlice(int(hashParams.SaltLength))
	if err != nil {
		return "", err
	}

	// Hash password.
	hash := argon2.IDKey([]byte(password), salt, uint32(hashParams.Iterations), hashParams.Memory, hashParams.Threads, uint32(hashParams.KeyLength))

	hashedPassword := Password{
		ArgonVersion: argon2.Version,
		Hash:         hash,
		Salt:         salt,
		Iterations:   hashParams.Iterations,
		Memory:       hashParams.Memory,
		Threads:      hashParams.Threads,
		KeyLength:    hashParams.KeyLength,
	}

	// Convert password struct to byte slice.
	hashedBytes, err := PasswordToBytes(&hashedPassword)
	if err != nil {
		return "", err
	}

	// Convert password byte slice into a string.
	return utils.EncodeString(hashedBytes), nil
}

// ComparePasswordAndHash compares a password to a hashed password to see if they are the same password.
func ComparePasswordAndHash(password, encodedHash string) (bool, error) {
	hashedBytes, err := utils.DecodeString(encodedHash)
	if err != nil {
		return false, err
	}

	pwd, err := PasswordFromBytes(hashedBytes)
	if err != nil {
		return false, err
	}

	if pwd.ArgonVersion != argon2.Version {
		return false, errors.New("mismatched argon versions between package and hash")
	}

	otherHash := argon2.IDKey([]byte(password), pwd.Salt, uint32(pwd.Iterations), pwd.Memory, pwd.Threads, uint32(pwd.KeyLength))

	if subtle.ConstantTimeCompare(pwd.Hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

// PasswordToBytes converts a Password struct into a slice of bytes.
func PasswordToBytes(hashedPassword *Password) ([]byte, error) {
	hashedBytes := new(bytes.Buffer)
	enc := gob.NewEncoder(hashedBytes)

	if err := enc.Encode(*hashedPassword); err != nil {
		return nil, err
	}

	return hashedBytes.Bytes(), nil
}

// PasswordFromBytes takes a slice of bytes and tries to decode them into a Password struct.
func PasswordFromBytes(hashedBytes []byte) (*Password, error) {
	hashedReader := bytes.NewBuffer(hashedBytes)
	dec := gob.NewDecoder(hashedReader)

	password := new(Password)

	err := dec.Decode(&password)
	if err != nil {
		return nil, err
	}

	return password, nil
}
