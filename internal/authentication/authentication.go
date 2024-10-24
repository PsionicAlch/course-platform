package authentication

import (
	"fmt"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication/errors"
	"github.com/PsionicAlch/psionicalch-home/internal/config"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/forms"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
)

type Authentication struct {
	utils.Loggers
	db            database.Database
	cookieWrapper *SecureCookieWrapper
}

func CreateAuthentication(db database.Database) (*Authentication, error) {
	loggers := utils.CreateLoggers("AUTHENTICATION")

	cookieWrapper, err := CreateSecureCookieWrapper()
	if err != nil {
		return nil, err
	}

	return &Authentication{
		Loggers:       loggers,
		db:            db,
		cookieWrapper: cookieWrapper,
	}, nil
}

func (auth *Authentication) SignUserIn(form *forms.SignUpForm, ipAddr string) (*http.Cookie, error) {
	// Check to make sure user doesn't already exist.
	user, err := auth.db.FindUserByEmail(form.Email)
	if err != nil {
		return nil, errors.CreateFailedToSignUserIn(err.Error())
	}

	if user != nil {
		return nil, errors.CreateUserAlreadyExists()
	}

	// Hash user's password.
	password, err := HashPassword(form.Password)
	if err != nil {
		return nil, errors.CreateFailedToSignUserIn(fmt.Sprintf("failed to hash user's password: %s", err))
	}

	// Add user to database.
	userId, err := auth.db.AddUser(form.Email, password)
	if err != nil {
		return nil, errors.CreateFailedToSignUserIn(err.Error())
	}

	// Generate authentication token.
	token, err := auth.db.CreateAuthenticationToken(userId, ipAddr)
	if err != nil {
		return nil, errors.CreateFailedToSignUserIn(err.Error())
	}

	// Save authentication token to secure cookie.
	encodedCookie, err := auth.cookieWrapper.Encode(token, form.RememberMe)
	if err != nil {
		return nil, errors.CreateFailedToSignUserIn(fmt.Sprintf("failed to create secure cookie: %s", err))
	}

	return encodedCookie, nil
}

func (auth *Authentication) GetUserFromCookie(cookies []*http.Cookie) (*models.UserModel, error) {
	for _, cookie := range cookies {
		if cookie.Name == config.GetWithoutError[string]("AUTH_COOKIE_NAME") {
			var authToken string
			if err := auth.cookieWrapper.Decode(cookie.Value, &authToken); err != nil {
				auth.ErrorLog.Println(err)
				return nil, err
			}

			auth.InfoLog.Println("Auth Token: ", authToken)
		}
	}

	return nil, nil
}
