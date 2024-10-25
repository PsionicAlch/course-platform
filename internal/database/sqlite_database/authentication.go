package sqlite_database

import (
	"database/sql"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/errors"
	"github.com/PsionicAlch/psionicalch-home/pkg/gatekeeper"
)

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
func (db *SQLiteDatabase) AddToken(token *gatekeeper.Token) error {
	query := `INSERT INTO tokens (id, token, token_type, valid_until, user_id, ip_address) VALUES (?, ?, ?, ?, ?, ?);`

	// Create token ID.
	id, err := database.GenerateID()
	if err != nil {
		return errors.CreateFailedToGenerateID(err.Error())
	}

	// Save token to the database.
	_, err = db.connection.Exec(query, id, token.Token, token.TokenType, token.ValidUntil, token.UserID, token.IPAddress)
	if err != nil {
		return errors.CreateFailedToCreateAuthenticationToken(err.Error())
	}

	return nil
}

func (db *SQLiteDatabase) GetToken(token string) (*gatekeeper.Token, error) {
	query := `SELECT token, token_type, valid_until, user_id, ip_address  FROM tokens WHERE token = ?;`
	row := db.connection.QueryRow(query, token)

	var dbToken string
	var tokenType string
	var validUntil time.Time
	var userId string
	var ipAddr string

	err := row.Scan(&dbToken, &tokenType, &validUntil, &userId, &ipAddr)
	if err != nil {
		return nil, err
	}

	return gatekeeper.NewToken(dbToken, tokenType, userId, ipAddr, validUntil)
}
