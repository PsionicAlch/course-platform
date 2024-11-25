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

func (db *SQLiteDatabase) GetAllAdminsPaginated(page, elements int) ([]*models.UserModel, error) {
	query := `SELECT id, name, surname, email, password, is_admin, is_author, affiliate_code, affiliate_points, created_at, updated_at FROM users WHERE is_admin = 1 LIMIT ? OFFSET ?;`

	offset := (page - 1) * elements

	var admins []*models.UserModel

	rows, err := db.connection.Query(query, elements, offset)
	if err != nil {
		db.ErrorLog.Printf("Failed to query database for all admin users (page %d): %s\n", page, err)
		return nil, err
	}

	for rows.Next() {
		var admin models.UserModel
		var isAdmin int
		var isAuthor int

		if err := rows.Scan(&admin.ID, &admin.Name, &admin.Surname, &admin.Email, &admin.Password, &isAdmin, &isAuthor, &admin.AffiliateCode, &admin.AffiliatePoints, &admin.CreatedAt, &admin.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from users table: %s\n", err)
			return nil, err
		}

		admin.IsAdmin = isAdmin == 1
		admin.IsAuthor = isAuthor == 1

		admins = append(admins, &admin)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to query database for all admin users (page %d): %s\n", page, err)
		return nil, err
	}

	return admins, nil
}

func (db *SQLiteDatabase) GetAllUsersPaginated(page, elements int) ([]*models.UserModel, error) {
	query := `SELECT id, name, surname, email, password, is_admin, is_author, affiliate_code, affiliate_points, created_at, updated_at FROM users WHERE is_admin = 0 AND is_author = 0 LIMIT ? OFFSET ?;`

	offset := (page - 1) * elements

	var users []*models.UserModel

	rows, err := db.connection.Query(query, elements, offset)
	if err != nil {
		db.ErrorLog.Printf("Failed to query database for all non-admin and non-author users (page %d): %s\n", page, err)
		return nil, err
	}

	for rows.Next() {
		var user models.UserModel
		var isAdmin int
		var isAuthor int

		if err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Password, &isAdmin, &isAuthor, &user.AffiliateCode, &user.AffiliatePoints, &user.CreatedAt, &user.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from users table: %s\n", err)
			return nil, err
		}

		user.IsAdmin = isAdmin == 1
		user.IsAuthor = isAuthor == 1

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to query database for all non-admin and non-author users (page %d): %s\n", page, err)
		return nil, err
	}

	return users, nil
}

func (db *SQLiteDatabase) GetAllAuthorsPaginated(page, elements int) ([]*models.UserModel, error) {
	query := `SELECT id, name, surname, email, password, is_admin, is_author, affiliate_code, affiliate_points, created_at, updated_at FROM users WHERE is_author = 1 LIMIT ? OFFSET ?;`

	offset := (page - 1) * elements

	var authors []*models.UserModel

	rows, err := db.connection.Query(query, elements, offset)
	if err != nil {
		db.ErrorLog.Printf("Failed to query database for all author users (page %d): %s\n", page, err)
		return nil, err
	}

	for rows.Next() {
		var author models.UserModel
		var isAdmin int
		var isAuthor int

		if err := rows.Scan(&author.ID, &author.Name, &author.Surname, &author.Password, &isAdmin, &isAuthor, &author.AffiliateCode, &author.AffiliatePoints, &author.CreatedAt, &author.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from users table: %s\n", err)
			return nil, err
		}

		author.IsAdmin = isAdmin == 1
		author.IsAuthor = isAuthor == 1

		authors = append(authors, &author)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to query database for all author users (page %d): %s\n", page, err)
		return nil, err
	}

	return authors, nil
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

func (db *SQLiteDatabase) GetUserByEmail(email string) (*models.UserModel, error) {
	query := `SELECT id, name, surname, email, password, is_admin, is_author, affiliate_code, affiliate_points, created_at, updated_at FROM users WHERE email = ?;`

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

func (db *SQLiteDatabase) GetUserByID(id string) (*models.UserModel, error) {
	query := `SELECT id, name, surname, email, password, is_admin, is_author, affiliate_code, affiliate_points, created_at, updated_at FROM users WHERE id = ?;`

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
