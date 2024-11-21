package internal

import (
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// AddToken adds a new token to the database. It works with transactions and single database connections.
// This function does database level uniqueness checks and will return a custom database.ErrTokenAlreadyExists
// which could be checked for on the caller side.
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
