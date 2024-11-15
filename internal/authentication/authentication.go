package authentication

import (
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/config"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
)

type Authentication struct {
	utils.Loggers
	Lifetime           time.Duration
	PasswordParameters *PasswordParameters
	CookiesManager     *CookieManager
	Database           database.Database
}

func SetupAuthentication(db database.Database) (*Authentication, error) {
	loggers := utils.CreateLoggers("AUTHENTICATION")

	lifetime := time.Duration(config.GetWithoutError[int]("AUTH_TOKEN_LIFETIME")) * time.Minute

	passwordParameters := DefaultPasswordParameters()

	cookiesManager, err := CreateCookieManager(lifetime)
	if err != nil {
		loggers.ErrorLog.Printf("Failed to create cookie manager: %s\n", err)
		return nil, err
	}

	auth := &Authentication{
		Loggers:            loggers,
		Lifetime:           lifetime,
		PasswordParameters: passwordParameters,
		CookiesManager:     cookiesManager,
		Database:           db,
	}

	return auth, nil
}

func (auth *Authentication) SignUserUp(name, surname, email, password, ipAddr string) (*http.Cookie, error) {
	// TODO: Switch to using database level uniqueness checks for whether the user exists or not.
	exists, err := auth.Database.UserExists(email)
	if err != nil {
		auth.ErrorLog.Printf("Failed to check if user with email \"%s\" already exists: %s\n", email, err)
		return nil, err
	}

	if exists {
		return nil, ErrUserExists
	}

	hashedPassword, err := auth.PasswordParameters.HashPassword(password)
	if err != nil {
		auth.ErrorLog.Printf("Failed to hash user's password: %s\n", err)
		return nil, err
	}

	userId, err := auth.Database.AddUser(name, surname, email, hashedPassword)
	if err != nil {
		auth.ErrorLog.Printf("Failed to save new user to the database: %s\n", err)
		return nil, err
	}

	token, err := NewToken()
	if err != nil {
		auth.ErrorLog.Printf("Failed to generate new token: %s\n", err)
		return nil, err
	}

	validUntil := time.Now().Add(auth.Lifetime)

	err = auth.Database.AddToken(token, AuthenticationToken, userId, ipAddr, validUntil)
	if err != nil {
		auth.ErrorLog.Printf("Failed to add %s token to the database: %s\n", AuthenticationToken, err)
		return nil, err
	}

	cookie, err := auth.CookiesManager.Encode(token)
	if err != nil {
		auth.ErrorLog.Printf("Failed to encode authentication cookie: %s\n", err)
		return nil, err
	}

	// TODO: Send email about new account creation.

	return cookie, nil
}

func (auth *Authentication) LogUserIn(email, password, ipAddr string) (*http.Cookie, error) {
	user, err := auth.Database.GetUser(email)
	if err != nil {
		auth.ErrorLog.Printf("Failed to find user (\"%s\") in database: %s\n", email, err)
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	match, err := ComparePasswordAndHash(password, user.Password)
	if err != nil {
		auth.ErrorLog.Printf("Failed to compare user (\"%s\") password: %s\n", email, err)
		return nil, err
	}

	if !match {
		return nil, ErrInvalidCredentials
	}

	token, err := NewToken()
	if err != nil {
		auth.ErrorLog.Printf("Failed to generate new authentication token for user (\"%s\"): %s\n", email, err)
		return nil, err
	}

	validUntil := time.Now().Add(auth.Lifetime)

	err = auth.Database.AddToken(token, AuthenticationToken, user.ID, ipAddr, validUntil)
	if err != nil {
		auth.ErrorLog.Printf("Failed to add %s token to the database: %s\n", AuthenticationToken, err)
		return nil, err
	}

	cookie, err := auth.CookiesManager.Encode(token)
	if err != nil {
		auth.ErrorLog.Printf("Failed to encode authentication cookie: %s\n", err)
		return nil, err
	}

	// TODO: Send email about new login just incase it wasn't the account holder who did it.

	return cookie, nil
}

func (auth *Authentication) LogUserOut(cookies []*http.Cookie) (*http.Cookie, error) {
	emptyCookie := auth.CookiesManager.EmptyCookie()

	for _, cookie := range cookies {
		if cookie.Name == auth.CookiesManager.CookieParams.Name {
			authToken, err := auth.CookiesManager.Decode(cookie.Value)
			if err != nil {
				auth.ErrorLog.Printf("Failed to decode auth cookie's value: %s\n", err)
				return emptyCookie, err
			}

			err = auth.Database.DeleteToken(authToken, AuthenticationToken)
			if err != nil {
				if err != database.ErrNoRowsAffected {
					auth.ErrorLog.Printf("Failed to delete authentication token: %s\n", err)
					return cookie, err
				}
			}

			// TODO: Reset session.

			return emptyCookie, nil
		}
	}

	return emptyCookie, nil
}

func (auth *Authentication) GetUserFromAuthCookie(cookies []*http.Cookie) (*models.UserModel, error) {
	for _, cookie := range cookies {
		if cookie.Name == auth.CookiesManager.CookieParams.Name {
			authToken, err := auth.CookiesManager.Decode(cookie.Value)
			if err != nil {
				auth.ErrorLog.Printf("Failed to decode auth cookie's value: %s\n", err)
				return nil, err
			}

			token, err := auth.Database.GetToken(authToken, AuthenticationToken)
			if err != nil {
				auth.ErrorLog.Printf("Failed to get authentication token from database: %s\n", err)
				return nil, err
			}

			valid := ValidateAuthenticationToken(token)
			if !valid {
				return nil, nil
			}

			user, err := auth.Database.GetUserByID(token.UserID)
			if err != nil {
				auth.ErrorLog.Printf("Failed to get user (\"%s\") from database: %s\n", token.UserID, err)
				return nil, err
			}

			return user, nil
		}
	}

	return nil, nil
}
