package internal

import (
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

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
