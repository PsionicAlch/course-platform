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

// TODO: The update logic is broken and needs to be fixed
func (db *SQLiteDatabase) BulkAddTutorials(tutorials []*models.TutorialModel) error {
	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("Failed to start transaction to bulk insert tutorials: %s\n", err)
		return err
	}

	for _, tutorial := range tutorials {
		// Try to insert the tutorial into the database.
		id, err := database.GenerateID()
		if err != nil {
			if err := tx.Rollback(); err != nil {
				db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
			}

			db.ErrorLog.Printf("Failed to generate ID for new tutorial: %s\n", err)
			return err
		}

		// Just in case we end up grabbing a new copy of the tutorial from the database we want to keep a local copy
		// of the new keywords.
		keywords := tutorial.Keywords

		if err := internal.AddTutorial(tx, id, tutorial.Title, tutorial.Slug, tutorial.Description, tutorial.ThumbnailURL, tutorial.BannerURL, tutorial.Content, tutorial.FileChecksum); err != nil {
			if err == database.ErrTutorialAlreadyExists {
				tutorial, err = internal.GetTutorialBySlug(tx, tutorial.Slug)
				if err != nil {
					if err := tx.Rollback(); err != nil {
						db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
					}

					db.ErrorLog.Printf("Failed to get tutorial from the database: %s\n", err)
					return err
				}
			} else {
				if err := tx.Rollback(); err != nil {
					db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
				}

				db.ErrorLog.Printf("Failed to add new tutorial to the database: %s\n", err)
				return err
			}
		}

		// In case this tutorial is getting updated we want to remove all keyword connections from the database.
		if err := internal.DeleteAllKeywordsFromTutorial(tx, tutorial.ID); err != nil {
			if err := tx.Rollback(); err != nil {
				db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
			}

			db.ErrorLog.Printf("Failed to remove keywords from tutorial \"%s\": %s\n", tutorial.Title, err)
			return err
		}

		for _, keyword := range keywords {
			// Generate a new ID for the new tutorials_keywords row.
			id, err := database.GenerateID()
			if err != nil {
				if err := tx.Rollback(); err != nil {
					db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
				}

				db.ErrorLog.Printf("Failed to generate ID for new tutorials_keywords row: %s\n", err)
				return err
			}

			// Assume the keyword doesn't exist yet so generate a new ID.
			keywordId, err := database.GenerateID()
			if err != nil {
				if err := tx.Rollback(); err != nil {
					db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
				}

				db.ErrorLog.Printf("Failed to generate ID for new keyword: %s\n", err)
				return err
			}

			keyword.ID = keywordId

			// Add the keyword to the tutorial via the pivot table. This function will work even if the keyword
			// already exists in the database.
			if err := internal.AddKeywordToTutorial(tx, id, tutorial.ID, keyword); err != nil {
				if err := tx.Rollback(); err != nil {
					db.ErrorLog.Printf("Failed to rollback transaction after an error occurred: %s\n", err)
				}

				db.ErrorLog.Printf("Failed to add new tutorials_keywords row: %s\n", err)
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
