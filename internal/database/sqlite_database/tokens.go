package sqlite_database

import (
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
)

func (db *SQLiteDatabase) AddToken(token, tokenType, userId, ipAddr string, validUntil time.Time) error {
	query := `INSERT INTO tokens (id, token, token_type, valid_until, user_id, ip_address) VALUES (?, ?, ?, ?, ?, ?);`

	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new ID for %s token: %s\n", tokenType, err)
		return err
	}

	result, err := db.connection.Exec(query, id, token, tokenType, validUntil, userId, ipAddr)
	if err != nil {
		db.ErrorLog.Printf("Failed to save %s token to the database: %s\n", tokenType, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to check rows affected after saving %s token to the database: %s\n", tokenType, err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Printf("No rows were affected by saving %s token to the database\n", tokenType)
		return database.ErrNoRowsAffected
	}

	return nil
}
