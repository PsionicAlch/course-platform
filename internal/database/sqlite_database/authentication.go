package sqlite_database

import (
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/config"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/errors"
)

const authTokenType = "authentication"

// AddUser adds the user to the database and returns their ID.
func (db *SQLiteDatabase) AddUser(email, password string) (string, error) {
	query := `INSERT INTO users (id, email, password) VALUES (?, ?, ?);`

	// Create user ID.
	id, err := database.GenerateID()
	if err != nil {
		return "", errors.CreateFailedToGenerateID(err.Error())
	}

	_, err = db.connection.Exec(query, id, email, password)
	if err != nil {
		return "", errors.CreateFailedToAddUserToDatabase(err.Error())
	}

	return id, nil
}

// CreateAuthenticationToken creates a new authentication token in the database and returns the token.
func (db *SQLiteDatabase) CreateAuthenticationToken(userId, ip string) (string, error) {
	query := `INSERT INTO tokens (id, token, token_type, valid_until, user_id, ip_address) VALUES (?, ?, ?, ?, ?, ?);`

	// Create token ID.
	id, err := database.GenerateID()
	if err != nil {
		return "", errors.CreateFailedToGenerateID(err.Error())
	}

	// Create token.
	token, err := database.GenerateToken()
	if err != nil {
		return "", errors.CreateFailedToGenerateToken(err.Error())
	}

	// Set date on which token expires.
	authTokenLifetime := time.Duration(config.GetWithoutError[int]("AUTH_TOKEN_LIFETIME"))
	validUntil := time.Now().Add(authTokenLifetime * time.Minute)

	// Save cookie to the database.
	_, err = db.connection.Exec(query, id, token, authTokenType, validUntil, userId, ip)
	if err != nil {
		return "", errors.CreateFailedToCreateAuthenticationToken(err.Error())
	}

	return token, nil
}
