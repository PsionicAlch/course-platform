package internal

import (
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
