package authentication

import (
	"net/http"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/session"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
)

type Authentication struct {
	utils.Loggers
	AuthenticationLifetime time.Duration
	PasswordResetLifetime  time.Duration
	PasswordParameters     *PasswordParameters
	CookiesManager         *CookieManager
	Database               database.Database
	Session                *session.Session
}

func SetupAuthentication(db database.Database, ses *session.Session, authLifetime, pwdResetLifetime time.Duration, cookieName, domainName, currentSecureCookieKey, previousSecureCookieKey string) (*Authentication, error) {
	loggers := utils.CreateLoggers("AUTHENTICATION")

	passwordParameters := DefaultPasswordParameters()

	cookiesManager, err := CreateCookieManager(authLifetime, cookieName, domainName, currentSecureCookieKey, previousSecureCookieKey)
	if err != nil {
		loggers.ErrorLog.Printf("Failed to create cookie manager: %s\n", err)
		return nil, err
	}

	auth := &Authentication{
		Loggers:                loggers,
		AuthenticationLifetime: authLifetime,
		PasswordResetLifetime:  pwdResetLifetime,
		PasswordParameters:     passwordParameters,
		CookiesManager:         cookiesManager,
		Database:               db,
		Session:                ses,
	}

	return auth, nil
}

func (auth *Authentication) SignUserUp(name, surname, email, password, ipAddr string) (*models.UserModel, *http.Cookie, error) {
	hashedPassword, err := auth.PasswordParameters.HashPassword(password)
	if err != nil {
		auth.ErrorLog.Printf("Failed to hash user's password: %s\n", err)
		return nil, nil, err
	}

	token, err := NewToken()
	if err != nil {
		auth.ErrorLog.Printf("Failed to generate new token: %s\n", err)
		return nil, nil, err
	}

	validUntil := time.Now().Add(auth.AuthenticationLifetime)

	user, err := auth.Database.AddNewUser(name, surname, email, hashedPassword, token, AuthenticationToken, ipAddr, validUntil)
	if err != nil {
		if err == database.ErrUserAlreadyExists {
			return nil, nil, ErrUserExists
		}

		auth.ErrorLog.Printf("Failed to save new user (\"%s\") to the database: %s\n", email, err)
		return nil, nil, err
	}

	cookie, err := auth.CookiesManager.Encode(token)
	if err != nil {
		auth.ErrorLog.Printf("Failed to encode authentication cookie: %s\n", err)
		return nil, nil, err
	}

	return user, cookie, nil
}

func (auth *Authentication) LogUserIn(email, password string) (*models.UserModel, *http.Cookie, error) {
	user, err := auth.Database.GetUserByEmail(email, database.All)
	if err != nil {
		auth.ErrorLog.Printf("Failed to find user (\"%s\") in database: %s\n", email, err)
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, ErrInvalidCredentials
	}

	match, err := ComparePasswordAndHash(password, user.Password)
	if err != nil {
		auth.ErrorLog.Printf("Failed to compare user (\"%s\") password: %s\n", email, err)
		return nil, nil, err
	}

	if !match {
		return nil, nil, ErrInvalidCredentials
	}

	token, err := NewToken()
	if err != nil {
		auth.ErrorLog.Printf("Failed to generate new authentication token for user (\"%s\"): %s\n", email, err)
		return nil, nil, err
	}

	validUntil := time.Now().Add(auth.AuthenticationLifetime)

	err = auth.Database.AddToken(token, AuthenticationToken, user.ID, validUntil)
	if err != nil {
		auth.ErrorLog.Printf("Failed to add %s token to the database: %s\n", AuthenticationToken, err)
		return nil, nil, err
	}

	cookie, err := auth.CookiesManager.Encode(token)
	if err != nil {
		auth.ErrorLog.Printf("Failed to encode authentication cookie: %s\n", err)
		return nil, nil, err
	}

	return user, cookie, nil
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

			valid := ValidateToken(token, AuthenticationToken)
			if !valid {
				auth.InfoLog.Printf("Authentication token was deemed invalid.")
				continue
			}

			user, err := auth.Database.GetUserByID(token.UserID, database.All)
			if err != nil {
				auth.ErrorLog.Printf("Failed to get user (\"%s\") from database: %s\n", token.UserID, err)
				return nil, err
			}

			return user, nil
		}
	}

	return nil, nil
}

func (auth *Authentication) GeneratePasswordResetToken(email string) (*models.UserModel, string, error) {
	user, err := auth.Database.GetUserByEmail(email, database.All)
	if err != nil {
		auth.ErrorLog.Printf("Failed to find user (\"%s\") in database: %s\n", email, err)
		return nil, "", err
	}

	if user == nil {
		return nil, "", ErrUnregisteredEmail
	}

	token, err := NewToken()
	if err != nil {
		auth.ErrorLog.Printf("Failed to generate new email token for user (\"%s\"): %s\n", email, err)
		return nil, "", err
	}

	validUntil := time.Now().Add(auth.PasswordResetLifetime)

	err = auth.Database.AddToken(token, EmailToken, user.ID, validUntil)
	if err != nil {
		auth.ErrorLog.Printf("Failed to add %s token to the database: %s\n", EmailToken, err)
		return nil, "", err
	}

	return user, token, nil
}

func (auth *Authentication) ValidateEmailToken(emailToken string) (bool, error) {
	token, err := auth.Database.GetToken(emailToken, EmailToken)
	if err != nil {
		auth.ErrorLog.Printf("Failed to get email token from database: %s\n", err)
		return false, err
	}

	return ValidateToken(token, EmailToken), nil
}

func (auth *Authentication) GetUserFromEmailToken(emailToken string) (*models.UserModel, error) {
	user, err := auth.Database.GetUserByToken(emailToken, EmailToken, database.All)
	if err != nil {
		auth.ErrorLog.Printf("Failed to get user using password reset token from database: %s\n", err)
		return nil, err
	}

	return user, nil
}

func (auth *Authentication) ChangeUserPassword(user *models.UserModel, password string) error {
	hashedPassword, err := auth.PasswordParameters.HashPassword(password)
	if err != nil {
		auth.ErrorLog.Printf("Failed to hash user's password: %s\n", err)
		return err
	}

	err = auth.Database.UpdateUserPassword(user.ID, hashedPassword)
	if err != nil {
		auth.ErrorLog.Printf("Failed to update user's password in the database: %s\n", err)
		return err
	}

	err = auth.Database.DeleteAllTokens(user.Email, AuthenticationToken)
	if err != nil {
		auth.ErrorLog.Printf("Failed to delete all of the user's (\"%s\") %s tokens: %s\n", user.Email, AuthenticationToken, err)
		return err
	}

	err = auth.Database.DeleteAllTokens(user.Email, EmailToken)
	if err != nil {
		auth.ErrorLog.Printf("Failed to delete all of the user's (\"%s\") %s tokens: %s\n", user.Email, EmailToken, err)
		return err
	}

	return nil
}

func (auth *Authentication) DeleteEmailToken(token string) error {
	err := auth.Database.DeleteToken(token, EmailToken)
	if err != nil {
		auth.ErrorLog.Printf("Failed to email token from the database: %s\n", err)
		return err
	}

	return nil
}

func (auth *Authentication) NewAdminUser(name, surname, email, password string) error {
	hashedPassword, err := auth.PasswordParameters.HashPassword(password)
	if err != nil {
		auth.ErrorLog.Printf("Failed to hash user's password: %s\n", err)
		return err
	}

	if err := auth.Database.NewAdminUser(name, surname, email, hashedPassword); err != nil {
		auth.ErrorLog.Printf("Failed to add new admin user: %s\n", err)
		return err
	}

	return nil
}

func (auth *Authentication) NewUser(name, surname, email, password string) error {
	hashedPassword, err := auth.PasswordParameters.HashPassword(password)
	if err != nil {
		auth.ErrorLog.Printf("Failed to hash user's password: %s\n", err)
		return err
	}

	if err := auth.Database.NewUser(name, surname, email, hashedPassword); err != nil {
		auth.ErrorLog.Printf("Failed to add new user: %s\n", err)
		return err
	}

	return nil
}
