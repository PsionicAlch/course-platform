package authentication

import (
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/config"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/gorilla/securecookie"
)

type SecureCookieWrapper struct {
	name             string
	domain           string
	expires          time.Time
	sameSite         http.SameSite
	secure           bool
	secureCookiesMap map[string]*securecookie.SecureCookie
}

func CreateSecureCookieWrapper() (*SecureCookieWrapper, error) {
	name := config.GetWithoutError[string]("AUTH_COOKIE_NAME")
	domain := config.GetWithoutError[string]("DOMAIN_NAME")
	expires := time.Now().Add(time.Duration(config.GetWithoutError[int]("AUTH_TOKEN_LIFETIME")) * time.Minute)
	sameSite := http.SameSiteLaxMode
	secure := config.InProduction()

	prevHashKey, prevBlockKey, err := DecodeKeys(config.GetWithoutError[string]("SECURE_COOKIE_PREVIOUS_HASH_KEY"), config.GetWithoutError[string]("SECURE_COOKIE_PREVIOUS_BLOCK_KEY"))
	if err != nil {
		return nil, err
	}

	hashKey, blockKey, err := DecodeKeys(config.GetWithoutError[string]("SECURE_COOKIE_HASH_KEY"), config.GetWithoutError[string]("SECURE_COOKIE_BLOCK_KEY"))
	if err != nil {
		return nil, err
	}

	secureCookiesMap := map[string]*securecookie.SecureCookie{
		"previous": securecookie.New(prevHashKey, prevBlockKey),
		"current":  securecookie.New(hashKey, blockKey),
	}

	cookieWrapper := &SecureCookieWrapper{
		name:             name,
		domain:           domain,
		expires:          expires,
		sameSite:         sameSite,
		secure:           secure,
		secureCookiesMap: secureCookiesMap,
	}

	return cookieWrapper, nil
}

func (wrapper *SecureCookieWrapper) Encode(value any, remember bool) (*http.Cookie, error) {
	encoded, err := securecookie.EncodeMulti(wrapper.name, value, wrapper.secureCookiesMap["current"])
	if err != nil {
		return nil, err
	}

	cookie := new(http.Cookie)
	cookie.Name = wrapper.name
	cookie.Value = encoded
	cookie.Path = "/"
	cookie.Domain = wrapper.domain

	if remember {
		cookie.Expires = wrapper.expires
	}

	cookie.HttpOnly = true
	cookie.SameSite = wrapper.sameSite
	cookie.Secure = wrapper.secure

	return cookie, nil
}

func (wrapper *SecureCookieWrapper) Decode(cookieValue string, dest any) error {
	return securecookie.DecodeMulti(wrapper.name, cookieValue, dest, wrapper.secureCookiesMap["current"], wrapper.secureCookiesMap["previous"])
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
