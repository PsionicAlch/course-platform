package internal

import (
	"database/sql"
	"errors"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"modernc.org/sqlite"
)

// FindKeyword is an internal function that can work with a normal database connection, a *sql.Stmt connection
// or even a *sql.Tx connection as long as they implement the functions that are required as specified by
// SqlDbFacade.
func FindKeyword(dbFacade SqlDbFacade, loggers utils.Loggers, keyword string) (*models.KeywordModel, error) {
	// Database query to get ID and keyword from keywords table.
	query := `SELECT id, keyword FROM keywords WHERE keyword = ?;`

	// Prepare the database connection to read in the information.
	row := dbFacade.QueryRow(query, keyword)

	// Create an intermediary variable to hold the associating data.
	keywordModel := new(models.KeywordModel)

	// Scan the database for the necessary information and handle errors. We don't care
	// about any errors related to not finding the data because we'll just send back
	// and empty struct pointer if nothing was found.
	err := row.Scan(&keywordModel.ID, &keywordModel.Keyword)
	if err != nil && err != sql.ErrNoRows {
		loggers.ErrorLog.Printf("Failed to add %s keyword to the keywords table: %s", keyword, err)
		return nil, err
	}

	return keywordModel, nil
}

// AddKeyword is an internal function that can work with a normal database connection, a *sql.Stmt connection
// or even a *sql.Tx connection as long as they implement the functions that are required as specified by
// SqlDbFacade.
func AddKeyword(dbFacade SqlDbFacade, loggers utils.Loggers, keyword string) (*models.KeywordModel, error) {
	// Database query to insert a new keyword into the keywords table..
	query := `INSERT INTO keywords (id, keyword) VALUES (?, ?);`

	// Create new ID.
	id, err := database.GenerateID()
	if err != nil {
		loggers.ErrorLog.Printf("Failed to generate new ID for %s keyword: %s", keyword, err)
		return nil, err
	}

	keywordModel := new(models.KeywordModel)

	// Try and insert the new row into the database table.
	result, err := dbFacade.Exec(query, id, keyword)
	if err != nil {
		// If the error is a unique constraint violation it means that the keyword already exists
		// so we'll just query the database for the necessary ID.
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == UNIQUE_CONSTRAINT_VIOLATION_ERROR_CODE {
			keywordModel, err := FindKeyword(dbFacade, loggers, keyword)
			if err != nil {
				loggers.ErrorLog.Printf("Failed to find %s keyword even though a copy should be in the database: %s", keyword, err)
			}

			return keywordModel, err
		}

		return nil, err
	}

	// Do some error checking to ensure that the new row has been added.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		loggers.ErrorLog.Printf("Failed to find out how many rows were affected whilst inserting keyword \"%s\" into the database: %s\n", keyword, err)
		return nil, err
	}

	if rowsAffected == 0 {
		loggers.ErrorLog.Printf("Failed to insert keyword \"%s\" into the database for unknown reason.", keyword)
		return nil, errors.New("failed to insert keyword into the database for unknown reason")
	}

	keywordModel.ID = id
	keywordModel.Keyword = keyword

	return keywordModel, nil
}

func AddKeywords(dbFacade SqlDbFacade, loggers utils.Loggers, keywords []string) ([]*models.KeywordModel, error) {
	var keywordModels []*models.KeywordModel

	// Range over each keyword and add them to the database.
	for _, keyword := range keywords {
		keywordModel, err := AddKeyword(dbFacade, loggers, keyword)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to insert new keyword \"%s\" into the database: %s", keyword, err)
			return nil, err
		}

		keywordModels = append(keywordModels, keywordModel)
	}

	return keywordModels, nil
}
