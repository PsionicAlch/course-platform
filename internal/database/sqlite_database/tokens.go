package sqlite_database

import (
	"database/sql"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
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

func (db *SQLiteDatabase) GetToken(token, tokenType string) (*models.TokenModel, error) {
	query := `SELECT id, token, token_type, valid_until, created_at, user_id, ip_address FROM tokens WHERE token = ? AND token_type = ?;`

	tokenStruct := new(models.TokenModel)

	row := db.connection.QueryRow(query, token, tokenType)

	if err := row.Scan(&tokenStruct.ID, &tokenStruct.Token, &tokenStruct.TokenType, &tokenStruct.ValidUntil, &tokenStruct.CreatedAt, &tokenStruct.UserID, &tokenStruct.IPAddr); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to get %s token from the database: %s\n", tokenType, err)
		return nil, err
	}

	return tokenStruct, nil
}

func (db *SQLiteDatabase) DeleteToken(token, tokenType string) error {
	query := `DELETE FROM tokens WHERE token = ? AND token_type = ?;`

	result, err := db.connection.Exec(query, token, tokenType)
	if err != nil {
		db.ErrorLog.Printf("Failed to delete token from database: %s\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to get rows affected after deleting token: %s\n", err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Println("No rows were affected after deleting token")
		return database.ErrNoRowsAffected
	}

	return nil
}
