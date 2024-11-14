package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
)

func (db *SQLiteDatabase) UserExists(email string) (bool, error) {
	query := `SELECT id FROM users WHERE email = ?;`
	row := db.connection.QueryRow(query, email)

	var id string
	err := row.Scan(&id)
	if err != nil {
		// User was not found so we can just return here.
		if err == sql.ErrNoRows {
			return false, nil
		}

		db.ErrorLog.Printf("Failed to check if user exists: %s\n", err.Error())
		return false, err
	}

	// Just in case Scan didn't return an ErrNoRows check if id is empty.
	return id == "", nil
}

func (db *SQLiteDatabase) AddUser(name, surname, email, password string) (string, error) {
	query := `INSERT INTO users (id, name, surname, email, password, affiliate_code) VALUES (?, ?, ?, ?, ?, ?);`

	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new ID for user (\"%s\"): %s\n", email, err.Error())
		return "", err
	}

	affiliate_code, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new affiliate code for user (\"%s\"): %s\n", email, err.Error())
		return "", err
	}

	result, err := db.connection.Exec(query, id, name, surname, email, password, affiliate_code)
	if err != nil {
		db.ErrorLog.Printf("Failed to save new user (\"%s\") to the database: %s\n", email, err.Error())
		return "", nil
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to see how many rows were affected after saving new user (\"%s\") to the database: %s\n", email, err)
		return "", err
	}

	if rowsAffected == 0 {
		return "", database.ErrNoRowsAffected
	}

	return id, nil
}
