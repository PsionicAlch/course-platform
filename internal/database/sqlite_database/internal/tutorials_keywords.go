package internal

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// AddKeywordToTutorial adds a new keyword-tutorial association in the database. This function works with either a
// database connection or a database transaction. This function wil NOT throw an error upon a unique constraint
// violation.
func AddKeywordToTutorial(dbFacade SqlDbFacade, id, keywordId, tutorialId string) error {
	query := `INSERT INTO tutorials_keywords (id, tutorial_id, keyword_id) VALUES (?, ?, ?);`

	result, err := dbFacade.Exec(query, id, tutorialId, keywordId)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return nil
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

// DeleteAllKeywordsFromTutorial removes all keyword associations from the tutorial with the provided ID. This function
// works with either a database connection or a database transaction.
func DeleteAllKeywordsFromTutorial(dbFacade SqlDbFacade, tutorialId string) error {
	query := `DELETE FROM tutorials_keywords WHERE tutorial_id = ?;`

	_, err := dbFacade.Exec(query, tutorialId)
	if err != nil {
		return err
	}

	return nil
}
