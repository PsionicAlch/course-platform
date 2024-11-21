package internal

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

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

func DeleteAllKeywordsFromTutorial(dbFacade SqlDbFacade, tutorialId string) error {
	query := `DELETE FROM tutorials_keywords WHERE tutorial_id = ?;`

	_, err := dbFacade.Exec(query, tutorialId)
	if err != nil {
		return err
	}

	return nil
}
