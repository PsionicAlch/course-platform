package internal

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

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

func GetKeywordByKeyword(dbFacade SqlDbFacade, keyword string) (*models.KeywordModel, error) {
	query := `SELECT id, keyword FROM keywords WHERE keyword = ?;`

	keywordModel := new(models.KeywordModel)

	row := dbFacade.QueryRow(query, keyword)
	if err := row.Scan(&keywordModel.ID, &keywordModel.Keyword); err != nil {
		return nil, err
	}

	return keywordModel, nil
}
