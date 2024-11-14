package authentication

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/config"
	"github.com/gorilla/securecookie"
)

type CookieParameters struct {
	Name     string
	Domain   string
	SameSite http.SameSite
	Secure   bool
	LifeTime time.Duration
}

type CookieManager struct {
	CookieParams         *CookieParameters
	CurrentSecureCookie  *securecookie.SecureCookie
	PreviousSecureCookie *securecookie.SecureCookie
}

func CreateCookieManager(lifetime time.Duration) (*CookieManager, error) {
	currentKeys, err := CreateSecureCookieKeys(config.GetWithoutError[string]("CURRENT_SECURE_COOKIE_KEY"))
	if err != nil {
		return nil, err
	}

	prevKeys, err := CreateSecureCookieKeys(config.GetWithoutError[string]("PREVIOUS_SECURE_COOKIE_KEY"))
	if err != nil {
		return nil, err
	}

	currentSecureCookie := securecookie.New(currentKeys.HashKey, currentKeys.BlockKey)
	previousSecureCookie := securecookie.New(prevKeys.HashKey, prevKeys.BlockKey)

	cookieParams := &CookieParameters{
		Name:     config.GetWithoutError[string]("AUTH_COOKIE_NAME"),
		Domain:   config.GetWithoutError[string]("DOMAIN_NAME"),
		SameSite: http.SameSiteLaxMode,
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
	cookie.Path = "/"
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
