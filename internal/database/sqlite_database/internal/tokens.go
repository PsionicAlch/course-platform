package internal

import (
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
)

func AddToken(dbFacade SqlDbFacade, id, token, tokenType, userId string, validUntil time.Time) error {
	query := `INSERT INTO tokens (id, token, token_type, valid_until, user_id) VALUES (?, ?, ?, ?, ?);`

	result, err := dbFacade.Exec(query, id, token, tokenType, validUntil, userId)
	if err != nil {
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
