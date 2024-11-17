package sqlite_database

import (
	"database/sql"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database/internal"
)

func (db *SQLiteDatabase) AddNewUser(name, surname, email, password, token, tokenType, ipAddr string, validUntil time.Time) (*models.UserModel, error) {
	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("Failed to begin transaction to add user (\"%s\") to the database: %s\n", email, err)
		return nil, err
	}

	userId, err := database.GenerateID()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback database changes: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to generate new ID for user (\"%s\"): %s\n", email, err)
		return nil, err
	}

	affiliateCode, err := database.GenerateID()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback database changes: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to generate new affiliate code for user (\"%s\"): %s\n", email, err)
		return nil, err
	}

	user, err := internal.AddUser(tx, userId, name, surname, email, password, affiliateCode)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback database changes: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to add user (\"%s\") to the database: %s\n", email, err)
		return nil, err
	}

	tokenId, err := database.GenerateID()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback database changes: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to generate new token ID for user (\"%s\"): %s\n", email, err)
		return nil, err
	}

	err = internal.AddToken(tx, tokenId, token, tokenType, user.ID, validUntil)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback database changes: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to save user's (\"%s\") %s token to the database: %s\n", email, tokenType, err)
		return nil, err
	}

	addressId, err := database.GenerateID()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback database changes: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to generate new IP address ID for user (\"%s\"): %s\n", email, err)
		return nil, err
	}

	err = internal.AddIPAddress(tx, addressId, user.ID, ipAddr)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback database changes: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to save IP address to the database: %s\n", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback database changes: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to commit transaction to the database when adding a new user (\"%s\"): %s\n", email, err)
		return nil, err
	}

	return user, nil
}

func (db *SQLiteDatabase) GetUser(email string) (*models.UserModel, error) {
	query := `SELECT id, name, surname, email, password, is_admin, is_author, affiliate_code, created_at, updated_at FROM users WHERE email = ?;`

	var isAdminInt int
	var isAuthorInt int
	user := new(models.UserModel)
	isAdmin := false
	isAuthor := false

	row := db.connection.QueryRow(query, email)
	if err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &isAdminInt, &isAuthorInt, &user.AffiliateCode, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			// Nothing was found so we can just send back nothing and handle it at the caller
			// end.
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to query the database for user (\"%s\"): %s\n", email, err)
		return nil, err
	}

	if isAdminInt == 1 {
		isAdmin = true
	}

	if isAuthorInt == 1 {
		isAuthor = true
	}

	user.IsAdmin = isAdmin
	user.IsAuthor = isAuthor

	return user, nil
}

func (db *SQLiteDatabase) GetUserByID(id string) (*models.UserModel, error) {
	query := `SELECT id, name, surname, email, password, is_admin, is_author, affiliate_code, created_at, updated_at FROM users WHERE id = ?;`

	var isAdminInt int
	var isAuthorInt int
	user := new(models.UserModel)
	isAdmin := false
	isAuthor := false

	row := db.connection.QueryRow(query, id)
	if err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &isAdminInt, &isAuthorInt, &user.AffiliateCode, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			// Nothing was found so we can just send back nothing and handle it at the caller
			// end.
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to query the database for user (\"%s\"): %s\n", id, err)
		return nil, err
	}

	if isAdminInt == 1 {
		isAdmin = true
	}

	if isAuthorInt == 1 {
		isAuthor = true
	}

	user.IsAdmin = isAdmin
	user.IsAuthor = isAuthor

	return user, nil
}

func (db *SQLiteDatabase) GetUserByToken(token, tokenType string) (*models.UserModel, error) {
	query := `SELECT users.id, users.name, users.surname, users.email, users.password, users.is_admin, users.is_author, users.affiliate_code, users.created_at, users.updated_at FROM tokens JOIN users ON tokens.user_id = users.id WHERE tokens.token = ? AND tokens.token_type = ? AND tokens.valid_until > CURRENT_TIMESTAMP;`

	var isAdminInt int
	var isAuthorInt int
	user := new(models.UserModel)
	isAdmin := false
	isAuthor := false

	row := db.connection.QueryRow(query, token, tokenType)
	if err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &isAdminInt, &isAuthorInt, &user.AffiliateCode, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			// Nothing was found so we can just send back nothing and handle it at the caller
			// end.
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to query the database for user using token (\"%s\"): %s\n", token, err)
		return nil, err
	}

	if isAdminInt == 1 {
		isAdmin = true
	}

	if isAuthorInt == 1 {
		isAuthor = true
	}

	user.IsAdmin = isAdmin
	user.IsAuthor = isAuthor

	return user, nil
}

func (db *SQLiteDatabase) UpdateUserPassword(userId, password string) error {
	query := `UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;`

	result, err := db.connection.Exec(query, password, userId)
	if err != nil {
		db.ErrorLog.Printf("Failed to update user's (\"%s\") password: %s\n", userId, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to get rows affected after updating user's (\"%s\") password: %s\n", userId, err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Printf("0 rows were affected after updating user's (\"%s\") password\n", userId)
		return database.ErrNoRowsAffected
	}

	return nil
}
