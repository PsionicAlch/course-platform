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
