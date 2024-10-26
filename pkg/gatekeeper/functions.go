package gatekeeper

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"net/http"
)

func newBytesSlice(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func newSalt(length int) ([]byte, error) {
	return newBytesSlice(length)
}

func passwordToBytes(password *GatekeeperPassword) ([]byte, error) {
	hashedBytes := new(bytes.Buffer)
	enc := gob.NewEncoder(hashedBytes)

	if err := enc.Encode(*password); err != nil {
		return nil, err
	}

	return hashedBytes.Bytes(), nil
}

func passwordFromBytes(password []byte) (*GatekeeperPassword, error) {
	hashedReader := bytes.NewBuffer(password)
	dec := gob.NewDecoder(hashedReader)

	passwordStruct := new(GatekeeperPassword)

	err := dec.Decode(passwordStruct)
	if err != nil {
		return nil, err
	}

	return passwordStruct, nil
}

func bytesToString(src []byte) string {
	return base64.RawStdEncoding.EncodeToString(src)
}

func stringToBytes(src string) ([]byte, error) {
	return base64.RawStdEncoding.Strict().DecodeString(src)
}

func newToken() (string, error) {
	tokenBytes, err := newBytesSlice(32)
	if err != nil {
		return "", err
	}

	return bytesToString(tokenBytes), nil
}

func redirect(w http.ResponseWriter, r *http.Request, url string, status ...int) {
	var statusCode int
	if len(status) > 0 {
		statusCode = status[0]
	} else {
		statusCode = http.StatusTemporaryRedirect
	}

	w.Header().Set("HX-Redirect", url)
	http.Redirect(w, r, url, statusCode)
}
