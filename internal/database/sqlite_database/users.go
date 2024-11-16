package sqlite_database

import (
	"database/sql"
	"time"

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
	return id != "", nil
}

func (db *SQLiteDatabase) AddUser(name, surname, email, password string) (*models.UserModel, error) {
	query := `INSERT INTO users (id, name, surname, email, password, affiliate_code, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new ID for user (\"%s\"): %s\n", email, err.Error())
		return nil, err
	}

	affiliate_code, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new affiliate code for user (\"%s\"): %s\n", email, err.Error())
		return nil, err
	}

	user := new(models.UserModel)
	user.ID = id
	user.Name = name
	user.Surname = surname
	user.Email = email
	user.Password = password
	user.AffiliateCode = affiliate_code
	user.IsAdmin = false
	user.IsAuthor = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := db.connection.Exec(query, user.ID, user.Name, user.Surname, user.Email, user.Password, user.AffiliateCode, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		db.ErrorLog.Printf("Failed to save new user (\"%s\") to the database: %s\n", email, err.Error())
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to see how many rows were affected after saving new user (\"%s\") to the database: %s\n", email, err)
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, database.ErrNoRowsAffected
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
