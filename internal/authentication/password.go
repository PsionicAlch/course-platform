package authentication

import (
	"bytes"
	"encoding/gob"
	"runtime"

	"golang.org/x/crypto/argon2"
)

type Password struct {
	ArgonVersion int
	Hash         []byte
	Salt         []byte
	Iterations   uint8
	Memory       uint32
	Threads      uint8
	KeyLength    uint8
}

type PasswordParameters struct {
	SaltLength uint8
	Iterations uint8
	Memory     uint32
	Threads    uint8
	KeyLength  uint8
}

func DefaultPasswordParameters() *PasswordParameters {
	return &PasswordParameters{
		SaltLength: 32,
		Iterations: uint8(runtime.NumCPU()),
		Memory:     64 * 1024,
		Threads:    uint8(runtime.NumCPU()),
		KeyLength:  32,
	}
}

func (params *PasswordParameters) HashPassword(password string) (string, error) {
	salt, err := NewSalt(int(params.SaltLength))
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, uint32(params.Iterations), params.Memory, params.Threads, uint32(params.KeyLength))

	passwordStruct := Password{
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
		return "", err
	}

	return BytesToString(passwordBytes), nil
}

func PasswordToBytes(password *Password) ([]byte, error) {
	hashedBytes := new(bytes.Buffer)
	enc := gob.NewEncoder(hashedBytes)

	if err := enc.Encode(*password); err != nil {
		return nil, err
	}

	return hashedBytes.Bytes(), nil
}

func NewSalt(length int) ([]byte, error) {
	return RandomBytes(uint(length))
}
