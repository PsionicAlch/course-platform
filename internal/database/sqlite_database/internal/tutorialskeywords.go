package internal

import (
	"errors"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"modernc.org/sqlite"
)

func AssociateKeywordWithTutorial(dbFacade SqlDbFacade, loggers utils.Loggers, tutorialID, keywordID string) error {
	query := `INSERT INTO tutorials_keywords (id, tutorial_id, keyword_id) VALUES (?, ?, ?);`

	id, err := database.GenerateID()
	if err != nil {
		loggers.ErrorLog.Printf("Failed to generate ID for tutorials_keywords new row: %s\n", err)
		return err
	}

	result, err := dbFacade.Exec(query, id, tutorialID, keywordID)
	if err != nil {
		// If the association is already there then just return out of the function.
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == UNIQUE_CONSTRAINT_VIOLATION_ERROR_CODE {
			return nil
		}

		loggers.ErrorLog.Printf("Failed to add row to tutorials_keywords table: %s\n", err)

		return err
	}

	// Do some error checking to ensure that the new row has been added.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		loggers.ErrorLog.Printf("Failed to find out how many rows were affected whilst inserting a new row into tutorials_keywords table: %s\n", err)
		return err
	}

	if rowsAffected == 0 {
		loggers.ErrorLog.Println("Failed to add row to tutorials_keywords table for unknown reason.")
		return errors.New("failed to add row to tutorials_keywords table for unknown reason")
	}

	return nil
}

func AssociateKeywordsWithTutorial(dbFacade SqlDbFacade, loggers utils.Loggers, tutorialId string, keywords []*models.KeywordModel) error {
	for _, keyword := range keywords {
		err := AssociateKeywordWithTutorial(dbFacade, loggers, tutorialId, keyword.ID)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to associate %s with the tutorial %s\n", keyword, tutorialId)
			return err
		}
	}

	return nil
}

func GetAllKeywordsForTutorial(dbFacade SqlDbFacade, loggers utils.Loggers, tutorialId string) ([]*models.KeywordModel, error) {
	query := `SELECT k.id, k.keyword FROM tutorials_keywords JOIN keywords k ON tk.keyword_id = k.id WHERE tk.tutorial_id = ?;`

	rows, err := dbFacade.Query(query, tutorialId)
	if err != nil {
		loggers.ErrorLog.Printf("Failed to read all the rows of the tutorials_keywords and keywords tables: %s\n", err)
		return nil, err
	}
	defer rows.Close()

	var keywords []*models.KeywordModel

	for rows.Next() {
		keyword := new(models.KeywordModel)

		err := rows.Scan(&keyword.ID, &keyword.Keyword)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to read row from tutorials_keywords table: %s\n", err)
			return nil, err
		}

		keywords = append(keywords, keyword)
	}

	if rows.Err() != nil {
		loggers.ErrorLog.Printf("Failed to read row from tutorials_keywords table: %s\n", err)
		return nil, err
	}

	return keywords, nil
}
