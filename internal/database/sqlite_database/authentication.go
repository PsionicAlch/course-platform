package sqlite_database

import (
	"database/sql"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/errors"
)

const authTokenType = "authentication"

// UserExists checks if a user with the given email address is in the database.
func (db *SQLiteDatabase) UserExists(email string) (bool, error) {
	query := `SELECT id FROM users WHERE email = ?;`
	row := db.connection.QueryRow(query, email)

	id := ""

	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

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
func (db *SQLiteDatabase) AddToken(token, tokenType string, validUntil time.Time, userId, ipAddr string) error {
	query := `INSERT INTO tokens (id, token, token_type, valid_until, user_id, ip_address) VALUES (?, ?, ?, ?, ?, ?);`

	// Create token ID.
	id, err := database.GenerateID()
	if err != nil {
		return errors.CreateFailedToGenerateID(err.Error())
	}

	// Save token to the database.
	_, err = db.connection.Exec(query, id, token, authTokenType, validUntil, userId, ipAddr)
	if err != nil {
		return errors.CreateFailedToCreateAuthenticationToken(err.Error())
	}

	return nil
}
