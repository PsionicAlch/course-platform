package internal

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// AddKeywords adds a new keyword row to the database. This function works with either a database connection or a
// database transaction. This function will NOT throw an error upon a unique constraint violation.
func AddKeyword(dbFacade SqlDbFacade, id, keyword string) error {
	query := `INSERT INTO keywords (id, keyword) VALUES (?, ?);`

	result, err := dbFacade.Exec(query, id, keyword)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return database.ErrKeywordAlreadyExists
		}

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

// GetKeywordByKeyword will retrieve a KeywordModel from the database using a given keyword. This function works with
// either a database connection or a database transaction.
func GetKeywordByKeyword(dbFacade SqlDbFacade, keyword string) (*models.KeywordModel, error) {
	query := `SELECT id, keyword FROM keywords WHERE keyword = ?;`

	keywordModel := new(models.KeywordModel)

	row := dbFacade.QueryRow(query, keyword)
	if err := row.Scan(&keywordModel.ID, &keywordModel.Keyword); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return keywordModel, nil
}
