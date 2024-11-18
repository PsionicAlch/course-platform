package authentication

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

type CookieParameters struct {
	Name     string
	Domain   string
	Path     string
	SameSite http.SameSite
	Secure   bool
	LifeTime time.Duration
}

type CookieManager struct {
	CookieParams         *CookieParameters
	CurrentSecureCookie  *securecookie.SecureCookie
	PreviousSecureCookie *securecookie.SecureCookie
}

func CreateCookieManager(lifetime time.Duration, cookieName, domainName, currentSecureCookieKey, previousSecureCookieKey string) (*CookieManager, error) {
	currentKeys, err := CreateSecureCookieKeys(currentSecureCookieKey)
	if err != nil {
		return nil, err
	}

	prevKeys, err := CreateSecureCookieKeys(previousSecureCookieKey)
	if err != nil {
		return nil, err
	}

	currentSecureCookie := securecookie.New(currentKeys.HashKey, currentKeys.BlockKey)
	previousSecureCookie := securecookie.New(prevKeys.HashKey, prevKeys.BlockKey)

	cookieParams := &CookieParameters{
		Name:     cookieName,
		Domain:   domainName,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		LifeTime: lifetime,
	}

	cookieManager := &CookieManager{
		CookieParams:         cookieParams,
		CurrentSecureCookie:  currentSecureCookie,
		PreviousSecureCookie: previousSecureCookie,
	}

	return cookieManager, nil
}

func (manager *CookieManager) Encode(token string) (*http.Cookie, error) {
	encoded, err := securecookie.EncodeMulti(manager.CookieParams.Name, token, manager.CurrentSecureCookie)
	if err != nil {
		return nil, fmt.Errorf("failed to encode cookie: %w", err)
	}

	cookie := new(http.Cookie)
	cookie.Name = manager.CookieParams.Name
	cookie.Value = encoded
	cookie.Path = manager.CookieParams.Path
	cookie.Domain = manager.CookieParams.Domain
	cookie.HttpOnly = true
	cookie.SameSite = manager.CookieParams.SameSite
	cookie.Secure = manager.CookieParams.Secure
	cookie.Expires = time.Now().Add(manager.CookieParams.LifeTime)

	return cookie, nil
}

func (manager *CookieManager) Decode(value string) (string, error) {
	var token string

	err := securecookie.DecodeMulti(manager.CookieParams.Name, value, &token, manager.CurrentSecureCookie, manager.PreviousSecureCookie)
	if err != nil {
		return "", fmt.Errorf("failed to decode cookie: %w", err)
	}

	return token, nil
}

func (manager *CookieManager) EmptyCookie() *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = manager.CookieParams.Name
	cookie.Value = ""
	cookie.Path = manager.CookieParams.Path
	cookie.Domain = manager.CookieParams.Domain
	cookie.HttpOnly = true
	cookie.SameSite = manager.CookieParams.SameSite
	cookie.Secure = manager.CookieParams.Secure
	cookie.Expires = time.Now().Add(-24 * time.Hour)

	return cookie
}
