package internal

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func AddKeywordToTutorial(dbFacade SqlDbFacade, id, tutorialId string, keyword *models.KeywordModel) error {
	query := `INSERT INTO tutorials_keywords (id, tutorial_id, keyword_id) VALUES (?, ?, ?);`

	if err := AddKeyword(dbFacade, keyword.ID, keyword.Keyword); err != nil {
		if err == database.ErrKeywordAlreadyExists {
			keyword, err = GetKeywordByKeyword(dbFacade, keyword.Keyword)
			if err != nil {
				return nil
			}
		} else {
			return nil
		}
	}

	result, err := dbFacade.Exec(query, id, tutorialId, keyword.ID)
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

func GetAllKeywordsForTutorial(dbFacade SqlDbFacade, tutorialId string) ([]*models.KeywordModel, error) {
	query := `SELECT k.id, k.keyword FROM tutorials_keywords tk JOIN keywords k ON tk.keyword_id = k.id WHERE tk.tutorial_id = ?;`

	var keywords []*models.KeywordModel

	rows, err := dbFacade.Query(query, tutorialId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var keyword models.KeywordModel

		if err := rows.Scan(&keyword.ID, &keyword.Keyword); err != nil {
			return nil, err
		}

		keywords = append(keywords, &keyword)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return keywords, nil
}
