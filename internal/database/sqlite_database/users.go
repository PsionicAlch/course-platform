package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

func (db *SQLiteDatabase) UserExists(email string) (bool, error) {
	query := `SELECT id FROM users WHERE email = ?;`
	row := db.connection.QueryRow(query, email)

	var id string
	if err := row.Scan(&id); err != nil {
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

func (db *SQLiteDatabase) GetUser(email string) (*models.UserModel, error) {
	query := `SELECT id, name, surname, email, password, affiliate_code, created_at, updated_at FROM users WHERE email = ?;`
	user := new(models.UserModel)

	row := db.connection.QueryRow(query, email)
	if err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.AffiliateCode, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			// Nothing was found so we can just send back nothing and handle it at the caller
			// end.
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to query the database for user (\"%s\"): %s\n", email, err)
		return nil, err
	}

	return user, nil
}
