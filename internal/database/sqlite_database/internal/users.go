package internal

import (
	"database/sql"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// AddUser adds a new user row to the database. This function will work with either a database connection or a database
// transaction.
func AddUser(dbFacade SqlDbFacade, id, name, surname, slug, email, password, affiliateCode string) (*models.UserModel, error) {
	query := `INSERT INTO users (id, name, surname, slug, email, password, affiliate_code, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	user := new(models.UserModel)
	user.ID = id
	user.Name = name
	user.Surname = surname
	user.Slug = slug
	user.Email = email
	user.Password = password
	user.AffiliateCode = affiliateCode
	user.IsAdmin = false
	user.IsAuthor = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.AffiliatePoints = 0

	result, err := dbFacade.Exec(query, user.ID, user.Name, user.Surname, user.Slug, user.Email, user.Password, user.AffiliateCode, user.CreatedAt, user.UpdatedAt)
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

// GetUserByID retrieves a UserModel from the database depending on the provided ID and authorization level of the user.
// This function will work with either a database connection or a database transaction.
func GetUserByID(dbFacade SqlDbFacade, id string, level database.AuthorizationLevel) (*models.UserModel, error) {
	query := `SELECT id, name, surname, slug, email, password, is_admin, is_author, affiliate_code, affiliate_points, created_at, updated_at FROM users WHERE id = ?`

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
	if err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Slug, &user.Email, &user.Password, &isAdmin, &isAuthor, &user.AffiliateCode, &user.AffiliatePoints, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	user.IsAdmin = isAdmin == 1
	user.IsAuthor = isAuthor == 1

	return &user, nil
}
