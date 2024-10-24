package gatekeeper

import (
	"net/http"
	"runtime"
	"time"
)

type Gatekeeper struct {
	lifeTime      time.Duration
	hashParams    *GatekeeperPasswordHashParameters
	cookieManager *GatekeeperCookieManager
	database      GatekeeperDatabase
}

// NewGatekeeper creates a new instance of the Gatekeeper authentication system with sensible defaults.
func NewGatekeeper(cookieName, websiteDomain string, authLifetime time.Duration, currentKeys, prevKeys *GatekeeperSecureCookieKeys, database GatekeeperDatabase) (*Gatekeeper, error) {
	hashParams := &GatekeeperPasswordHashParameters{
		SaltLength: 32,
		Iterations: uint8(runtime.NumCPU()),
		Memory:     64 * 1024,
		Threads:    uint8(runtime.NumCPU()),
		KeyLength:  32,
	}

	cookieParams := &CookieParameters{
		Name:         cookieName,
		Domain:       websiteDomain,
		SameSite:     http.SameSiteLaxMode,
		Secure:       true,
		Lifetime:     authLifetime,
		CurrentKeys:  currentKeys,
		PreviousKeys: prevKeys,
	}

	cookieManager := CreateCookieManager(cookieParams)

	gt := &Gatekeeper{
		lifeTime:      authLifetime,
		hashParams:    hashParams,
		cookieManager: cookieManager,
		database:      database,
	}

	return gt, nil
}

// SignUserIn handles the logic around checking if the user already exists, hashing their password, saving the user to the database,
// creating a new authentication token, saving the token to the database, and creating a new cookie with all the information
// required to authenticate the user later on.
func (gatekeeper *Gatekeeper) SignUserIn(email, password, ipAddr string, rememberMe bool) (*http.Cookie, error) {
	// Try and fetch a user with that email address.
	userExists, err := gatekeeper.database.UserExists(email)
	if err != nil {
		return nil, createFailedToFindUserByEmail(email, err.Error())
	}

	// Make sure that the returned ID string is empty.
	if userExists {
		return nil, createUserAlreadyExists(email)
	}

	// Hash user's password.
	hashedPassword, err := gatekeeper.hashParams.HashPassword(password)
	if err != nil {
		return nil, createFailedToHashPassword(err.Error())
	}

	// Add user to the database.
	userId, err := gatekeeper.database.AddUser(email, hashedPassword)
	if err != nil {
		return nil, createFailedToAddUserToDatabase(err.Error())
	}

	// Generate new authentication token.
	token, err := NewToken()
	if err != nil {
		return nil, createFailedToGenerateNewToken(err.Error())
	}

	// Set the token type to be "authentication".
	tokenType := "authentication"

	// Set the expiry time and date for the token based on the authentication token's lifetime.
	validUntil := time.Now().Add(gatekeeper.lifeTime)

	// Save the token in the database.
	err = gatekeeper.database.AddToken(token, tokenType, validUntil, userId, ipAddr)
	if err != nil {
		return nil, createFailedToAddTokenToDatabase(err.Error())
	}

	// Create a new authentication cookie to help authenticate the user later on.
	encodedCookie, err := gatekeeper.cookieManager.Encode(token, rememberMe)
	if err != nil {
		return nil, createFailedToCreateAuthenticationCookie(err.Error())
	}

	return encodedCookie, nil
}
