package gatekeeper

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
)

type GatekeeperSecureCookieKeys struct {
	hashKey  []byte
	blockKey []byte
}

func CreateGatekeeperSecureCookieKeys(key string) (*GatekeeperSecureCookieKeys, error) {
	if key == "" {
		return new(GatekeeperSecureCookieKeys), nil
	}

	hashKey, blockKey, found := strings.Cut(key, "$")
	if !found {
		return nil, createInvalidGatekeeperKey("")
	}

	hashKeyBytes, err := stringToBytes(hashKey)
	if err != nil {
		return nil, createInvalidGatekeeperKey(err.Error())
	}

	blockKeyBytes, err := stringToBytes(blockKey)
	if err != nil {
		return nil, createInvalidGatekeeperKey(err.Error())
	}

	secureCookieKeys := &GatekeeperSecureCookieKeys{
		hashKey:  hashKeyBytes,
		blockKey: blockKeyBytes,
	}

	return secureCookieKeys, nil
}

func GenerateGatekeeperKey() (string, error) {
	hashKey, err := generateHashKey()
	if err != nil {
		return "", err
	}

	blockKey, err := generateBlockKey()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s$%s", hashKey, blockKey), nil
}

func generateHashKey() (string, error) {
	byteSlice, err := newBytesSlice(64)
	if err != nil {
		return "", createFailedToGenerateSecureCookieKey("hash", err.Error())
	}

	return bytesToString(byteSlice), nil
}

func generateBlockKey() (string, error) {
	byteSlice, err := newBytesSlice(32)
	if err != nil {
		return "", createFailedToGenerateSecureCookieKey("block", err.Error())
	}

	return bytesToString(byteSlice), nil
}

type CookieParameters struct {
	name         string
	domain       string
	sameSite     http.SameSite
	secure       bool
	lifetime     time.Duration
	currentKeys  *GatekeeperSecureCookieKeys
	previousKeys *GatekeeperSecureCookieKeys
}

type GatekeeperCookieManager struct {
	parameters           *CookieParameters
	currentSecureCookie  *securecookie.SecureCookie
	previousSecureCookie *securecookie.SecureCookie
}

func CreateCookieManager(params *CookieParameters) *GatekeeperCookieManager {
	currentSecureCookie := securecookie.New(params.currentKeys.hashKey, params.currentKeys.blockKey)
	previousSecureCookie := securecookie.New(params.previousKeys.hashKey, params.previousKeys.blockKey)

	cookieManager := &GatekeeperCookieManager{
		parameters:           params,
		currentSecureCookie:  currentSecureCookie,
		previousSecureCookie: previousSecureCookie,
	}

	return cookieManager
}

func (manager *GatekeeperCookieManager) Encode(value any, remember bool) (*http.Cookie, error) {
	encoded, err := securecookie.EncodeMulti(manager.parameters.name, value, manager.currentSecureCookie)
	if err != nil {
		return nil, createFailedToEncodeSecureCookie(err.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = manager.parameters.name
	cookie.Value = encoded
	cookie.Path = "/"
	cookie.Domain = manager.parameters.domain
	cookie.HttpOnly = true
	cookie.SameSite = manager.parameters.sameSite
	cookie.Secure = manager.parameters.secure
	if remember {
		cookie.Expires = time.Now().Add(manager.parameters.lifetime)
	}

	return cookie, nil
}

func (manager *GatekeeperCookieManager) Decode(value string, dest any) error {
	return securecookie.DecodeMulti(manager.parameters.name, value, dest, manager.currentSecureCookie, manager.previousSecureCookie)
}
