package internal

import (
	"database/sql"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// AddUser is an internal function that allows you to add a new user to the database with either a standalone connection
// or through a database connection. This function will return a user model or an error. It will fail on unique constraint
// violation with a custom database.ErrUserAlreadyExists. This can be used to do thread safe and runtime safe data uniqueness
// tests instead of doing two separate database calls where you check if the user exists first and then try to add them.
// Those kinds of tests can fail because a user could have been added after you checked but before you added your user.
func AddUser(dbFacade SqlDbFacade, id, name, surname, email, password, affiliateCode string) (*models.UserModel, error) {
	query := `INSERT INTO users (id, name, surname, email, password, affiliate_code, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

	user := new(models.UserModel)
	user.ID = id
	user.Name = name
	user.Surname = surname
	user.Email = email
	user.Password = password
	user.AffiliateCode = affiliateCode
	user.IsAdmin = false
	user.IsAuthor = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.AffiliatePoints = 0

	result, err := dbFacade.Exec(query, user.ID, user.Name, user.Surname, user.Email, user.Password, user.AffiliateCode, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return nil, database.ErrUserAlreadyExists
		}

		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, database.ErrNoRowsAffected
	}

	return user, nil
}

func GetUserByID(dbFacade SqlDbFacade, id string, level database.AuthorizationLevel) (*models.UserModel, error) {
	query := `SELECT id, name, surname, email, password, is_admin, is_author, affiliate_code, affiliate_points, created_at, updated_at FROM users WHERE id = ?`

	switch level {
	case database.User:
		query += ` AND is_admin = 0 AND is_author = 0`
	case database.Admin:
		query += ` AND is_admin = 1`
	case database.Author:
		query += ` AND is_author = 1`
	}

	var isAdmin int
	var isAuthor int
	var user models.UserModel

	row := dbFacade.QueryRow(query, id)
	if err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &user.Password, &isAdmin, &isAuthor, &user.AffiliateCode, &user.AffiliatePoints, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	user.IsAdmin = isAdmin == 1
	user.IsAuthor = isAuthor == 1

	return &user, nil
}

func UpdateAffiliatePoints(dbFacade SqlDbFacade, userId string, affiliatePoints uint) error {
	query := `UPDATE users SET affiliate_points = ? WHERE id = ?;`

	result, err := dbFacade.Exec(query, affiliatePoints, userId)
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
