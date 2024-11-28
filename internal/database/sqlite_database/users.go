package sqlite_database

import (
	"database/sql"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database/internal"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func (db *SQLiteDatabase) GetUsers(term string, level database.AuthorizationLevel) ([]*models.UserModel, error) {
	query := `SELECT id, name, surname, email, password, affiliate_code, affiliate_points, is_admin, is_author, created_at, updated_at FROM users WHERE (LOWER(id) LIKE '%' || ? || '%' OR LOWER(name) LIKE '%' || ? || '%' OR LOWER(surname) LIKE '%' || ? || '%' OR LOWER(email) LIKE '%' || ? || '%' OR LOWER(affiliate_code) LIKE '%' || ? || '%')`

	args := []any{term, term, term, term, term}

	switch level {
	case database.User:
		query += ` AND is_admin = 0 AND is_author = 0`
	case database.Admin:
		query += ` AND is_admin = 1`
	case database.Author:
		query += ` AND is_author = 1`
	}

	var users []*models.UserModel

	rows, err := db.connection.Query(query, args...)
	if err != nil {
		db.ErrorLog.Printf("Failed to query database for all users according to the search term \"%s\" and authorization level \"%s\": %s\n", term, level, err)
		return nil, err
	}

	for rows.Next() {
		var user models.UserModel
		var isAdmin int
		var isAuthor int

		if err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.AffiliateCode, &user.AffiliatePoints, &isAdmin, &isAuthor, &user.CreatedAt, &user.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from users table: %s\n", err)
			return nil, err
		}

		user.IsAdmin = isAdmin == 1
		user.IsAuthor = isAuthor == 1

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to query database for all users according to the search term \"%s\" and authorization level \"%s\": %s\n", term, level, err)
		return nil, err
	}

	return users, nil
}

func (db *SQLiteDatabase) GetUsersPaginated(term string, level database.AuthorizationLevel, page, elements uint) ([]*models.UserModel, error) {
	query := `SELECT id, name, surname, email, password, affiliate_code, affiliate_points, is_admin, is_author, created_at, updated_at FROM users WHERE (LOWER(id) LIKE '%' || ? || '%' OR LOWER(name) LIKE '%' || ? || '%' OR LOWER(surname) LIKE '%' || ? || '%' OR LOWER(email) LIKE '%' || ? || '%' OR LOWER(affiliate_code) LIKE '%' || ? || '%')`

	offset := (page - 1) * elements

	args := []any{term, term, term, term, term}

	switch level {
	case database.User:
		query += ` AND is_admin = 0 AND is_author = 0`
	case database.Admin:
		query += ` AND is_admin = 1`
	case database.Author:
		query += ` AND is_author = 1`
	}

	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?;`
	args = append(args, elements, offset)

	var users []*models.UserModel

	rows, err := db.connection.Query(query, args...)
	if err != nil {
		db.ErrorLog.Printf("Failed to query database for all users (page %d) according to the search term \"%s\" and authorization level \"%s\": %s\n", page, term, level, err)
		return nil, err
	}

	for rows.Next() {
		var user models.UserModel
		var isAdmin int
		var isAuthor int

		if err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &user.AffiliateCode, &user.AffiliatePoints, &isAdmin, &isAuthor, &user.CreatedAt, &user.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from users table: %s\n", err)
			return nil, err
		}

		user.IsAdmin = isAdmin == 1
		user.IsAuthor = isAuthor == 1

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to query database for all users (page %d) according to the search term \"%s\" and authorization level \"%s\": %s\n", page, term, level, err)
		return nil, err
	}

	return users, nil
}

// The only function that I believe has a viable reason to do transactions.
func (db *SQLiteDatabase) AddNewUser(name, surname, email, password, token, tokenType, ipAddr string, validUntil time.Time) (*models.UserModel, error) {
	// Generate all the IDs required for the database transaction first. I don't want the transaction to fail
	// because I couldn't generate an ID.
	userId, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new ID for user (\"%s\"): %s\n", email, err)
		return nil, err
	}

	affiliateCode, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new affiliate code for user (\"%s\"): %s\n", email, err)
		return nil, err
	}

	tokenId, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new token ID for user (\"%s\"): %s\n", email, err)
		return nil, err
	}

	addressId, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new IP address ID for user (\"%s\"): %s\n", email, err)
		return nil, err
	}

	// With all the setup out the way, start a new transaction. The only thing that can fail now is database calls.
	// The reason for the transaction is that I don't want the user saved to the database if the token or the IP
	// address couldn't have been saved. I'd rather the user start again than have my database contain fractured data.
	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("Failed to begin transaction to add user (\"%s\") to the database: %s\n", email, err)
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

	err = internal.AddToken(tx, tokenId, token, tokenType, user.ID, validUntil)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback database changes: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to save user's (\"%s\") %s token to the database: %s\n", email, tokenType, err)
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

	// Now that all the database calls have been successfully made we can commit the changes as one big chunk. Hopefully
	// committing the changes won't fail but if it does the user can just resubmit the signup form.
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback database changes: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to commit transaction to the database when adding a new user (\"%s\"): %s\n", email, err)
		return nil, err
	}

	return user, nil
}

func (db *SQLiteDatabase) NewUser(name, surname, email, password string) error {
	query := `INSERT INTO users (id, name, surname, email, password, affiliate_code) VALUES (?, ?, ?, ?, ?, ?);`

	userId, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new ID for user (\"%s\"): %s\n", email, err)
		return err
	}

	affiliateCode, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new affiliate code for user (\"%s\"): %s\n", email, err)
		return err
	}

	result, err := db.connection.Exec(query, userId, name, surname, email, password, affiliateCode)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return database.ErrUserAlreadyExists
		}

		db.ErrorLog.Printf("Failed to add new user to the database: %s\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to query the database for the rows affected after adding new user: %s\n", err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Println("0 rows were affected after adding new user to the database")
		return database.ErrNoRowsAffected
	}

	return nil
}

// AddAdminUser adds a new admin user to the database. This function should not be made available to the application because
// there should never be a reason for the actual application to use this function.
func (db *SQLiteDatabase) NewAdminUser(name, surname, email, password string) error {
	query := `INSERT INTO users (id, name, surname, email, password, affiliate_code, is_admin) VALUES (?, ?, ?, ?, ?, ?, 1);`

	userId, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new ID for user (\"%s\"): %s\n", email, err)
		return err
	}

	affiliateCode, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new affiliate code for user (\"%s\"): %s\n", email, err)
		return err
	}

	result, err := db.connection.Exec(query, userId, name, surname, email, password, affiliateCode)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return database.ErrUserAlreadyExists
		}

		db.ErrorLog.Printf("Failed to add new admin user to the database: %s\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to query the database for the rows affected after adding new admin user: %s\n", err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Println("0 rows were affected after adding new admin user to the database")
		return database.ErrNoRowsAffected
	}

	return nil
}

func (db *SQLiteDatabase) GetUserByEmail(email string, level database.AuthorizationLevel) (*models.UserModel, error) {
	query := `SELECT id, name, surname, email, password, is_admin, is_author, affiliate_code, affiliate_points, created_at, updated_at FROM users WHERE email = ?`

	switch level {
	case database.User:
		query += ` AND is_admin = 0 AND is_author = 0`
	case database.Admin:
		query += ` AND is_admin = 1`
	case database.Author:
		query += ` AND is_author = 1`
	}

	var isAdminInt int
	var isAuthorInt int
	user := new(models.UserModel)
	isAdmin := false
	isAuthor := false

	row := db.connection.QueryRow(query, email)
	if err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &isAdminInt, &isAuthorInt, &user.AffiliateCode, &user.AffiliatePoints, &user.CreatedAt, &user.UpdatedAt); err != nil {
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

func (db *SQLiteDatabase) GetUserByID(id string, level database.AuthorizationLevel) (*models.UserModel, error) {
	query := `SELECT id, name, surname, email, password, is_admin, is_author, affiliate_code, affiliate_points, created_at, updated_at FROM users WHERE id = ?`

	switch level {
	case database.User:
		query += ` AND is_admin = 0 AND is_author = 0`
	case database.Admin:
		query += ` AND is_admin = 1`
	case database.Author:
		query += ` AND is_author = 1`
	}

	var isAdminInt int
	var isAuthorInt int
	user := new(models.UserModel)
	isAdmin := false
	isAuthor := false

	row := db.connection.QueryRow(query, id)
	if err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &isAdminInt, &isAuthorInt, &user.AffiliateCode, &user.AffiliatePoints, &user.CreatedAt, &user.UpdatedAt); err != nil {
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
	query := `SELECT users.id, users.name, users.surname, users.email, users.password, users.is_admin, users.is_author, users.affiliate_code, users.affiliate_points, users.created_at, users.updated_at FROM tokens JOIN users ON tokens.user_id = users.id WHERE tokens.token = ? AND tokens.token_type = ? AND tokens.valid_until > CURRENT_TIMESTAMP;`

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

func (db *SQLiteDatabase) CountUsers() (uint, error) {
	query := `SELECT COUNT(id) FROM users;`

	var count uint

	row := db.connection.QueryRow(query)
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		db.ErrorLog.Printf("Failed to count the number of users in the database: %s\n", err)
		return 0, err
	}

	return count, nil
}

func (db *SQLiteDatabase) AddAuthorStatus(userId string) error {
	query := `UPDATE users SET is_author = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?;`

	_, err := db.connection.Exec(query, userId)
	if err != nil {
		db.ErrorLog.Printf("Failed to update user's (\"%s\") author status: %s\n", userId, err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) RemoveAuthorStatus(userId string) error {
	query := `UPDATE users SET is_author = 0, updated_at = CURRENT_TIMESTAMP WHERE id = ?;`

	_, err := db.connection.Exec(query, userId)
	if err != nil {
		db.ErrorLog.Printf("Failed to update user's (\"%s\") author status: %s\n", userId, err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) AddAdminStatus(userId string) error {
	query := `UPDATE users SET is_admin = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?;`

	_, err := db.connection.Exec(query, userId)
	if err != nil {
		db.ErrorLog.Printf("Failed to update user's (\"%s\") admin status: %s\n", userId, err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) RemoveAdminStatus(userId string) error {
	query := `UPDATE users SET is_admin = 0, updated_at = CURRENT_TIMESTAMP WHERE id = ?;`

	_, err := db.connection.Exec(query, userId)
	if err != nil {
		db.ErrorLog.Printf("Failed to update user's (\"%s\") admin status: %s\n", userId, err)
		return err
	}

	return nil
}
