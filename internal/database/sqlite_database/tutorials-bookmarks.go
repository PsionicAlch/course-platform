package sqlite_database

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database/internal"
)

func (db *SQLiteDatabase) UserBookmarkedTutorial(userId, slug string) (bool, error) {
	bookmarked, err := internal.UserBookmarkedTutorial(db.connection, userId, slug)
	if err != nil {
		db.ErrorLog.Printf("Failed to see if user bookmarked tutorial %s in the database: %s\n", slug, err)
		return false, err
	}

	return bookmarked, nil
}

func (db *SQLiteDatabase) UserBookmarkTutorial(userId, slug string) error {
	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate ID for new tutorials_bookmarks database row: %s\n", err)
		return err
	}

	if err := internal.UserBookmarkTutorial(db.connection, id, userId, slug); err != nil {
		db.ErrorLog.Printf("Failed to insert new row into tutorials_bookmarks database table: %s\n", err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) UserUnbookmarkTutorial(userId, slug string) error {
	if err := internal.UserUnbookmarkTutorial(db.connection, userId, slug); err != nil {
		db.ErrorLog.Printf("Failed to delete row from tutorials_bookmarks database table: %s\n", err)
		return err
	}

	return nil
}
