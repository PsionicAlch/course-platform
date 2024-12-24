package sqlite_database

import (
	"database/sql"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database/internal"
)

func (db *SQLiteDatabase) AddToken(token, tokenType, userId string, validUntil time.Time) error {
	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new ID for %s token: %s\n", tokenType, err)
		return err
	}

	err = internal.AddToken(db.connection, id, token, tokenType, userId, validUntil)
	if err != nil {
		db.ErrorLog.Printf("Failed to save %s token to the database: %s\n", tokenType, err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) GetToken(token, tokenType string) (*models.TokenModel, error) {
	query := `SELECT id, token, token_type, valid_until, created_at, user_id FROM tokens WHERE token = ? AND token_type = ?;`

	tokenStruct := new(models.TokenModel)

	row := db.connection.QueryRow(query, token, tokenType)

	if err := row.Scan(&tokenStruct.ID, &tokenStruct.Token, &tokenStruct.TokenType, &tokenStruct.ValidUntil, &tokenStruct.CreatedAt, &tokenStruct.UserID); err != nil {
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

	_, err := db.connection.Exec(query, token, tokenType)
	if err != nil {
		db.ErrorLog.Printf("Failed to delete token from database: %s\n", err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) DeleteAllTokens(email, tokenType string) error {
	query := `DELETE FROM tokens WHERE token_type = ? AND user_id IN (SELECT id FROM users WHERE email = ?);`

	_, err := db.connection.Exec(query, tokenType, email)
	if err != nil {
		db.ErrorLog.Printf("Failed to delete all user's (\"%s\") %s tokens from database: %s\n", email, tokenType, err)
		return err
	}

	return nil
}
