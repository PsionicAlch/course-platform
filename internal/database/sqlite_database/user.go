package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

func (db *SQLiteDatabase) FindUserByEmail(email string) (*models.UserModel, error) {
	query := `SELECT id, email, created_at, updated_at FROM users WHERE email = ?;`
	row := db.connection.QueryRow(query, email)

	user := new(models.UserModel)

	err := row.Scan(&user.ID, &user.Email, &user.Created_At, &user.Update_At)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}
