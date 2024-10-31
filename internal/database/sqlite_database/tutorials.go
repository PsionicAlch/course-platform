package sqlite_database

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database/internal"
)

// FindKeyword queries the keywords database table for the associate keyword. It can work with a normal
// database connection, a *sql.Stmt connection or even a *sql.Tx connection as long as they implement
// the functions that are required as specified by DBFacade.
func (db *SQLiteDatabase) FindKeyword(keyword string) (*models.KeywordModel, error) {
	keywordModel, err := internal.FindKeyword(db.connection, db.Loggers, keyword)
	if err != nil {
		db.ErrorLog.Printf("Failed to query keywords table for %s keyword: %s", keyword, err)
		return nil, err
	}

	return keywordModel, nil
}

// AddKeyword will insert the new keyword into the keywords table and return a KeywordModel or an error.
// It can work with a normal database connection, a *sql.Stmt connection or even a *sql.Tx connection as
// long as they implement the functions that are required as specified by DBFacade.
func (db *SQLiteDatabase) AddKeyword(keyword string) (*models.KeywordModel, error) {
	keywordModel, err := internal.AddKeyword(db.connection, db.Loggers, keyword)
	if err != nil {
		db.ErrorLog.Printf("Failed to add %s keyword to the keywords table: %s", keyword, err)
		return nil, err
	}

	return keywordModel, nil
}

// AddKeywordBulk will add a slice of keywords to the database, making sure to check whether they're already in the database
// first. (Really wish there was more errors, that way I could just try and insert it and have the database handle the
// unique index errors and then I can just check if that error was raised...)
func (db *SQLiteDatabase) AddKeywordBulk(keywords []string) ([]*models.KeywordModel, error) {
	// Start a transaction.
	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("Failed to start a new database transaction to bulk insert keywords: %s\n", err)
		return nil, err
	}

	keywordModels, err := internal.AddKeywords(tx, db.Loggers, keywords)
	if err != nil {
		db.ErrorLog.Printf("Failed to add keywords to database: %s\n", err)
		return nil, err
	}

	// Commit the transaction if all operations succeed
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		db.ErrorLog.Printf("Failed to commit keyword bulk insert transaction: %s", err)
		return nil, err
	}

	return keywordModels, nil
}

func (db *SQLiteDatabase) GetAllTutorials() ([]*models.TutorialModel, error) {
	tutorials, err := internal.GetAllTutorials(db.connection, db.Loggers)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all tutorials from the database: %s\n", err)
		return nil, err
	}

	return tutorials, nil
}

func (db *SQLiteDatabase) GetTutorialBySlug(slug string) (*models.TutorialModel, error) {
	tutorialModel, err := internal.GetTutorialBySlug(db.connection, db.Loggers, slug)
	if err != nil {
		db.ErrorLog.Printf("Failed to find tutorial by slug: %s\n", err)
		return nil, err
	}

	return tutorialModel, nil
}

func (db *SQLiteDatabase) AddNewTutorial(title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum string, keywords []string) error {
	// Begin a transaction since this function will result in a lot of SELECT and INSERT statements.
	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("Failed to begin transaction to insert a new tutorial: %s\n", err)
		return err
	}

	err = internal.AddTutorial(tx, db.Loggers, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, keywords)
	if err != nil {
		db.ErrorLog.Printf("Failed to add tutorial to the database: %s", err)
		return err
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.ErrorLog.Printf("Failed to commit transaction for inserting a new tutorial to the database: %s\n", err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) AddNewTutorialBulk(tutorials []*models.TutorialModel) error {
	// Begin a transaction.
	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("Failed to begin transaction for bulk insert into tutorials table: %s\n", err)
		return err
	}

	// Add tutorials to the database.
	err = internal.AddTutorials(tx, db.Loggers, tutorials)
	if err != nil {
		db.ErrorLog.Printf("Failed to add tutorials to database: %s\n", err)
		return err
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.ErrorLog.Printf("Failed to commit transaction for bulk insert into tutorials table: %s\n", err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) UpdateTutorial(id string, tutorial *models.TutorialModel) error {
	// Begin a transaction.
	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("Failed to begin transaction for updating tutorial: %s\n", err)
		return err
	}

	// Update tutorial.
	err = internal.UpdateTutorial(tx, db.Loggers, id, tutorial)
	if err != nil {
		db.ErrorLog.Printf("Failed to update tutorial: %s\n", err)
		return err
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.ErrorLog.Printf("Failed to commit transaction for updating tutorial: %s\n", err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) UpdateTutorialBulk(tutorials []*models.TutorialModel) error {
	// Start a transaction.
	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("failed to begin transaction for bulk update to tutorials table: %s\n", err)
		return err
	}

	// Update tutorials.
	err = internal.UpdateTutorials(tx, db.Loggers, tutorials)
	if err != nil {
		db.ErrorLog.Printf("Failed to update all tutorials: %s\n", err)
		return err
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.ErrorLog.Printf("Failed to commit bulk update transaction: %s\n", err)
		return err
	}

	return nil
}
