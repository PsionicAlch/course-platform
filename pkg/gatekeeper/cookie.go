package gatekeeper

import (
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/gorilla/securecookie"
)

type GatekeeperSecureCookieKeys struct {
	HashKey  []byte
	BlockKey []byte
}

func CreateGatekeeperSecureCookieKeys(hashKey, blockKey string) (*GatekeeperSecureCookieKeys, error) {
	hashKeyBytes, err := StringToBytes(hashKey)
	if err != nil {
		// TODO: Create dedicated error.
		return nil, err
	}

	blockKeyBytes, err := StringToBytes(blockKey)
	if err != nil {
		// TODO: Create dedicated error.
		return nil, err
	}

	secureCookieKeys := &GatekeeperSecureCookieKeys{
		HashKey:  hashKeyBytes,
		BlockKey: blockKeyBytes,
	}

	return secureCookieKeys, nil
}

func GenerateHashKey() (string, error) {
	byteSlice, err := NewBytesSlice(64)
	if err != nil {
		// TODO: Create dedicated error.
		return "", err
	}

	return BytesToString(byteSlice), nil
}

func GenerateBlockKey() (string, error) {
	byteSlice, err := NewBytesSlice(32)
	if err != nil {
		// TODO: Create dedicated error.
		return "", err
	}

	return BytesToString(byteSlice), nil
}

func DecodeKeys(encodedHashKey, encodedBlockKey string) ([]byte, []byte, error) {
	decodedHashKey, err := utils.DecodeString(encodedHashKey)
	if err != nil {
		return nil, nil, err
	}

	decodedBlockKey, err := utils.DecodeString(encodedBlockKey)
	if err != nil {
		return nil, nil, err
	}

	return decodedHashKey, decodedBlockKey, nil
}

type CookieParameters struct {
	Name         string
	Domain       string
	SameSite     http.SameSite
	Secure       bool
	Lifetime     time.Duration
	CurrentKeys  *GatekeeperSecureCookieKeys
	PreviousKeys *GatekeeperSecureCookieKeys
}

type GatekeeperCookieManager struct {
	Parameters           *CookieParameters
	CurrentSecureCookie  *securecookie.SecureCookie
	PreviousSecureCookie *securecookie.SecureCookie
}

func CreateCookieManager(params *CookieParameters) *GatekeeperCookieManager {
	currentSecureCookie := securecookie.New(params.CurrentKeys.HashKey, params.CurrentKeys.BlockKey)
	previousSecureCookie := securecookie.New(params.PreviousKeys.HashKey, params.PreviousKeys.BlockKey)

	cookieManager := &GatekeeperCookieManager{
		Parameters:           params,
		CurrentSecureCookie:  currentSecureCookie,
		PreviousSecureCookie: previousSecureCookie,
	}

	return cookieManager
}

func (manager *GatekeeperCookieManager) Encode(value any, remember bool) (*http.Cookie, error) {
	encoded, err := securecookie.EncodeMulti(manager.Parameters.Name, value, manager.CurrentSecureCookie)
	if err != nil {
		// TODO: Create dedicated error.
		return nil, err
	}

	cookie := new(http.Cookie)
	cookie.Name = manager.Parameters.Name
	cookie.Value = encoded
	cookie.Path = "/"
	cookie.Domain = manager.Parameters.Domain

	if remember {
		cookie.Expires = time.Now().Add(manager.Parameters.Lifetime)
	}

	cookie.HttpOnly = true
	cookie.SameSite = manager.Parameters.SameSite
	cookie.Secure = manager.Parameters.Secure

	return cookie, nil
}

func (manager *GatekeeperCookieManager) Decode(value string, dest any) error {
	return securecookie.DecodeMulti(manager.Parameters.Name, value, dest, manager.CurrentSecureCookie, manager.PreviousSecureCookie)
}
