package sqlite_database

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database/internal"
)

func (db *SQLiteDatabase) GetAllTutorials() ([]*models.TutorialModel, error) {
	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("Failed to start a transaction to query database for all tutorials: %s\n", err)
		return nil, err
	}

	tutorials, err := internal.GetAllTutorials(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to query database for all tutorials: %s\n", err)
		return nil, err
	}

	for _, tutorial := range tutorials {
		keywords, err := internal.GetAllKeywordsForTutorial(tx, tutorial.ID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
			}

			db.ErrorLog.Printf("Failed to query database for all keywords related to \"%s\": %s\n", tutorial.Title, err)
			return nil, err
		}

		tutorial.Keywords = keywords
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to commit transaction after query database for all tutorials: %s\n", err)
		return nil, err
	}

	return tutorials, nil
}

func (db *SQLiteDatabase) GetAllTutorialsPaginated(page, elements int) ([]*models.TutorialModel, error) {
	tutorials, err := internal.GetAllTutorialsPaginated(db.connection, page, elements)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all tutorials (paginated) from the database: %s\n", err)
		return nil, err
	}

	for _, tutorial := range tutorials {
		keywords, err := internal.GetAllKeywordsForTutorial(db.connection, tutorial.ID)
		if err != nil {
			db.ErrorLog.Printf("Failed to get all keywords for \"%s\" tutorial from the database: %s\n", tutorial.Title, err)
			return nil, err
		}

		tutorial.Keywords = keywords
	}

	return tutorials, nil
}

func (db *SQLiteDatabase) SearchTutorialsPaginated(term string, page, elements int) ([]*models.TutorialModel, error) {
	tutorials, err := internal.SearchTutorialsPaginated(db.connection, term, page, elements)
	if err != nil {
		db.ErrorLog.Printf("Failed to search for tutorials (paginated) from the database: %s\n", err)
		return nil, err
	}

	for _, tutorial := range tutorials {
		keywords, err := internal.GetAllKeywordsForTutorial(db.connection, tutorial.ID)
		if err != nil {
			db.ErrorLog.Printf("Failed to get all keywords for \"%s\" tutorial from the database: %s\n", tutorial.Title, err)
			return nil, err
		}

		tutorial.Keywords = keywords
	}

	return tutorials, nil
}

func (db *SQLiteDatabase) GetTutorialByID(id string) (*models.TutorialModel, error) {
	tutorial, err := internal.GetTutorialByID(db.connection, id)
	if err != nil {
		db.ErrorLog.Printf("Failed to get tutorial by id (\"%s\") from the database: %s\n", id, err)
		return nil, err
	}

	return tutorial, nil
}

func (db *SQLiteDatabase) GetTutorialBySlug(slug string) (*models.TutorialModel, error) {
	tutorial, err := internal.GetTutorialBySlug(db.connection, slug)
	if err != nil {
		db.ErrorLog.Printf("Failed to get tutorial by slug (\"%s\") from the database: %s\n", slug, err)
		return nil, err
	}

	keywords, err := internal.GetAllKeywordsForTutorial(db.connection, tutorial.ID)
	if err != nil {
		db.ErrorLog.Printf("Failed to get keywords for tutorial by slug (\"%s\") from the database: %s\n", slug, err)
		return nil, err
	}

	tutorial.Keywords = keywords

	return tutorial, nil
}

func (db *SQLiteDatabase) BulkAddTutorials(tutorials []*models.TutorialModel) error {
	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("Failed to start transaction to bulk insert tutorials: %s\n", err)
		return err
	}

	for _, tutorial := range tutorials {
		tutorialId, err := database.GenerateID()
		if err != nil {
			if err := tx.Rollback(); err != nil {
				db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
			}

			db.ErrorLog.Printf("Failed generate ID for new tutorial: %s\n", err)
			return err
		}

		if err := internal.AddTutorial(tx, tutorialId, tutorial.Title, tutorial.Slug, tutorial.Description, tutorial.ThumbnailURL, tutorial.BannerURL, tutorial.Content, tutorial.FileChecksum, tutorial.FileKey); err != nil {
			if err := tx.Rollback(); err != nil {
				db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
			}

			db.ErrorLog.Printf("Failed insert new tutorial in the database: %s\n", err)
			return err
		}

		for _, keyword := range tutorial.Keywords {
			id, err := database.GenerateID()
			if err != nil {
				if err := tx.Rollback(); err != nil {
					db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
				}

				db.ErrorLog.Printf("Failed generate ID new tutorials_keywords row: %s\n", err)
				return err
			}

			keywordId, err := database.GenerateID()
			if err != nil {
				if err := tx.Rollback(); err != nil {
					db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
				}

				db.ErrorLog.Printf("Failed generate ID for (possibly) new keyword: %s\n", err)
				return err
			}

			keyword.ID = keywordId

			if err := internal.AddKeywordToTutorial(tx, id, tutorialId, keyword); err != nil {
				if err := tx.Rollback(); err != nil {
					db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
				}

				db.ErrorLog.Printf("Failed to add keyword to tutorials_keywords table: %s\n", err)
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed commit transaction after bulk inserting tutorials: %s\n", err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) BulkUpdateTutorials(tutorials []*models.TutorialModel) error {
	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("Failed to start transaction to bulk update tutorials: %s\n", err)
		return err
	}

	for _, tutorial := range tutorials {
		if err := internal.UpdateTutorial(tx, tutorial.ID, tutorial.Title, tutorial.Slug, tutorial.Description, tutorial.ThumbnailURL, tutorial.BannerURL, tutorial.Content, tutorial.FileChecksum, tutorial.FileKey, tutorial.AuthorID); err != nil {
			if err := tx.Rollback(); err != nil {
				db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
			}

			db.ErrorLog.Printf("Failed to update tutorial in the database: %s\n", err)
			return err
		}

		if err := internal.DeleteAllKeywordsFromTutorial(tx, tutorial.ID); err != nil {
			if err := tx.Rollback(); err != nil {
				db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
			}

			db.ErrorLog.Printf("Failed to remove all keywords from tutorial in tutorials_keywords table: %s\n", err)
			return err
		}

		for _, keyword := range tutorial.Keywords {
			id, err := database.GenerateID()
			if err != nil {
				if err := tx.Rollback(); err != nil {
					db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
				}

				db.ErrorLog.Printf("Failed generate ID new tutorials_keywords row: %s\n", err)
				return err
			}

			keywordId, err := database.GenerateID()
			if err != nil {
				if err := tx.Rollback(); err != nil {
					db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
				}

				db.ErrorLog.Printf("Failed generate ID for (possibly) new keyword: %s\n", err)
				return err
			}

			keyword.ID = keywordId

			if err := internal.AddKeywordToTutorial(tx, id, tutorial.ID, keyword); err != nil {
				if err := tx.Rollback(); err != nil {
					db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
				}

				db.ErrorLog.Printf("Failed to add keyword to tutorials_keywords table: %s\n", err)
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed commit transaction after bulk updating tutorials: %s\n", err)
		return err
	}

	return nil
}
