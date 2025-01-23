package internal

import (
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// AddToken adds a new token to the database. This function works with either a database connection or a database
// transaction. This function WILL throw a ErrTokenAlreadyExists error upon a unique constraint violation.
func AddToken(dbFacade SqlDbFacade, id, token, tokenType, userId string, validUntil time.Time) error {
	query := `INSERT INTO tokens (id, token, token_type, valid_until, user_id) VALUES (?, ?, ?, ?, ?);`

	result, err := dbFacade.Exec(query, id, token, tokenType, validUntil, userId)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return database.ErrTokenAlreadyExists
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNoRowsAffected
	}

	return nil
}
