package gatekeeper

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
)

func NewBytesSlice(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func NewSalt(length int) ([]byte, error) {
	return NewBytesSlice(length)
}

func PasswordToBytes(password *GatekeeperPassword) ([]byte, error) {
	hashedBytes := new(bytes.Buffer)
	enc := gob.NewEncoder(hashedBytes)

	if err := enc.Encode(*password); err != nil {
		return nil, err
	}

	return hashedBytes.Bytes(), nil
}

func PasswordFromBytes(password []byte) (*GatekeeperPassword, error) {
	hashedReader := bytes.NewBuffer(password)
	dec := gob.NewDecoder(hashedReader)

	passwordStruct := new(GatekeeperPassword)

	err := dec.Decode(&password)
	if err != nil {
		return nil, err
	}

	return passwordStruct, nil
}

func BytesToString(src []byte) string {
	return base64.RawStdEncoding.EncodeToString(src)
}

func StringToBytes(src string) ([]byte, error) {
	return base64.RawStdEncoding.Strict().DecodeString(src)
}

func NewToken() (string, error) {
	tokenBytes, err := NewBytesSlice(32)
	if err != nil {
		return "", err
	}

	return BytesToString(tokenBytes), nil
}
