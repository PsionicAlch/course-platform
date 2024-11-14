package authentication

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomBytes(length uint) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func BytesToString(src []byte) string {
	return base64.RawStdEncoding.EncodeToString(src)
}

func StringToBytes(src string) ([]byte, error) {
	return base64.RawStdEncoding.Strict().DecodeString(src)
}
