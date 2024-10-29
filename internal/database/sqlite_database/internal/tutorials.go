package internal

import (
	"database/sql"
	"fmt"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"modernc.org/sqlite"
)

func GetAllTutorials(dbFacade SqlDbFacade, loggers utils.Loggers) ([]*models.TutorialModel, error) {
	// SQL query to get all tutorials from tutorials table.
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, created_at, updated_at FROM tutorials ORDER BY updated_at DESC, title ASC;`

	// Query the database.
	rows, err := dbFacade.Query(query)
	if err != nil {
		loggers.ErrorLog.Printf("Failed to read all the rows of the tutorials table: %s\n", err)
		return nil, err
	}
	defer rows.Close()

	var tutorials []*models.TutorialModel

	for rows.Next() {
		// Some more intermediate variables to hold each row.
		tutorial := new(models.TutorialModel)
		// SQLite doesn't have booleans so we have to use integers as an intermediate type.
		var publishedInt int

		// Read the data from the table row into the intermediary variables.
		err := rows.Scan(
			&tutorial.ID,
			&tutorial.Title,
			&tutorial.Slug,
			&tutorial.Description,
			&tutorial.ThumbnailURL,
			&tutorial.BannerURL,
			&tutorial.Content,
			&publishedInt,
			&tutorial.AuthorID,
			&tutorial.FileChecksum,
			&tutorial.CreatedAt,
			&tutorial.UpdatedAt,
		)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to query row of tutorials table: %s\n", err)
			return nil, err
		}

		// Convert the published column to a boolean.
		tutorial.Published = publishedInt == 1

		// Get all associated keywords.
		keywords, err := GetAllKeywordsForTutorial(dbFacade, loggers, tutorial.ID)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to get all keywords associated with tutorial %s: %s\n", tutorial.ID, err)
		}

		tutorial.Keywords = keywords

		// Add the individual tutorial to the slice of tutorials.
		tutorials = append(tutorials, tutorial)
	}

	// Check for any errors after iterating through all the rows.
	if rows.Err() != nil {
		loggers.ErrorLog.Printf("Found an error after querying all the tutorials table's rows: %s\n", err)
		return nil, err
	}

	return tutorials, nil
}

func FindTutorialBySlugWithoutKeywords(dbFacade SqlDbFacade, loggers utils.Loggers, slug string) (*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, created_at, updated_at FROM tutorials WHERE slug = ?;`

	tutorialModel := new(models.TutorialModel)
	var publishedInt int

	row := dbFacade.QueryRow(query, slug)
	err := row.Scan(
		&tutorialModel.ID,
		&tutorialModel.Title,
		&tutorialModel.Slug,
		&tutorialModel.Description,
		&tutorialModel.ThumbnailURL,
		&tutorialModel.BannerURL,
		&tutorialModel.Content,
		&publishedInt,
		&tutorialModel.AuthorID,
		&tutorialModel.FileChecksum,
		&tutorialModel.CreatedAt,
		&tutorialModel.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return tutorialModel, nil
		}

		loggers.ErrorLog.Printf("Failed to query row of tutorials table: %s\n", err)
		return nil, err
	}

	// Convert the published column to a boolean.
	tutorialModel.Published = publishedInt == 1

	return tutorialModel, nil
}

func AddTutorial(dbFacade SqlDbFacade, loggers utils.Loggers, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum string, keywords []string) error {
	query := `INSERT INTO tutorials (id, title, slug, description, thumbnail_url, banner_url, content, file_checksum) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

	// Generate new ULID based ID for the tutorial.
	id, err := database.GenerateID()
	if err != nil {
		loggers.ErrorLog.Printf("Failed to generate new ID for tutorial: %s\n", err)
		return err
	}

	// Execute the SQL query and pass in the required variables.
	_, err = dbFacade.Exec(query, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum)
	if err != nil {
		// Since you can't have multiple tutorials with the same slug, if this tutorial has already been added to the
		// database we'll get a unique constraint violation error from SQLite. So just in case this function has errored
		// out in the past after the tutorial was added to the database but before the keywords were added and associated
		// let's try and get the tutorial from the database. We'll use the slug since that has to be unique.
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == UNIQUE_CONSTRAINT_VIOLATION_ERROR_CODE {
			tutorialModel, err := FindTutorialBySlugWithoutKeywords(dbFacade, loggers, slug)
			if err != nil {
				loggers.ErrorLog.Printf("%s is already in the database but failed to retrieve it: %s\n", title, err)
				return nil
			}

			id = tutorialModel.ID

			// Probably shouldn't do this but since I don't want to do the select statement unless the tutorial is already
			// in the database (because it seems like a "good premature optimization" to make) but once we have the tutorial's
			// ID we still need to get back to the rest of the application and this was the only sane way I could see to
			// get back on track with the rest of the function. Older me will probably be wise enough to solve this.
			goto ContinueWithKeywords
		}

		loggers.ErrorLog.Printf("Failed to save tutorial \"%s\" to the database: %s\n", title, err)
		return err
	}

ContinueWithKeywords:

	// Add keywords to the database.
	keywordModels, err := AddKeywords(dbFacade, loggers, keywords)
	if err != nil {
		loggers.ErrorLog.Printf("Failed to save keywords for %s: %s\n", title, err)
		return nil
	}

	err = AssociateKeywordsWithTutorial(dbFacade, loggers, id, keywordModels)
	if err != nil {
		loggers.ErrorLog.Printf("Failed to associate keywords with %s: %s\n", title, err)
	}

	return nil
}

func AddTutorials(dbFacade SqlDbFacade, loggers utils.Loggers, tutorials []*models.TutorialModel) error {
	// Loop through each tutorial and insert it.
	for _, tutorial := range tutorials {
		var keywords []string
		for _, keyword := range tutorial.Keywords {
			keywords = append(keywords, keyword.Keyword)
		}

		err := AddTutorial(dbFacade, loggers, tutorial.Title, tutorial.Slug, tutorial.Description, tutorial.ThumbnailURL, tutorial.BannerURL, tutorial.Content, tutorial.FileChecksum, keywords)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to add tutorial %s to the database: %s", tutorial.Title, err)
			return err
		}
	}

	return nil
}

func UpdateTutorial(dbFacade SqlDbFacade, loggers utils.Loggers, tutorialId string, tutorial *models.TutorialModel) error {
	query := `UPDATE tutorials SET title = ?, slug = ?, description = ?, thumbnail_url = ?, banner_url = ?, content = ?, published = 0, author_id = null, file_checksum = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;`

	// Update tutorial.
	result, err := dbFacade.Exec(query, tutorial.Title, tutorial.Slug, tutorial.Description, tutorial.ThumbnailURL, tutorial.BannerURL, tutorial.Content, tutorial.FileChecksum, tutorialId)
	if err != nil {
		loggers.ErrorLog.Printf("Failed update tutorial %s: %s\n", tutorialId, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		loggers.ErrorLog.Printf("Failed to get the affected rows after updating tutorial %s: %s\n", tutorialId, err)
		return err
	}

	if rowsAffected == 0 {
		loggers.ErrorLog.Printf("Failed to update tutorial %s. 0 rows were affected by this update\n", tutorialId)
		return fmt.Errorf("failed to update tutorial %s. 0 rows were affected by this update", tutorialId)
	}

	// Update keywords.
	for _, keyword := range tutorial.Keywords {
		keywordModel, err := AddKeyword(dbFacade, loggers, keyword.Keyword)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to add keyword to the database whilst updating %s: %s\n", tutorial.Title, err)
			return err
		}

		err = AssociateKeywordWithTutorial(dbFacade, loggers, tutorialId, keywordModel.ID)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to associate %s with %s whilst updating tutorial: %s", keywordModel.Keyword, tutorial.Title, err)
			return err
		}
	}

	return nil
}

func UpdateTutorials(dbFacade SqlDbFacade, loggers utils.Loggers, tutorials []*models.TutorialModel) error {
	for _, tutorial := range tutorials {
		err := UpdateTutorial(dbFacade, loggers, tutorial.ID, tutorial)
		if err != nil {
			loggers.ErrorLog.Printf("Failed to update tutorial: %s\n", err)
			return err
		}
	}

	return nil
}
