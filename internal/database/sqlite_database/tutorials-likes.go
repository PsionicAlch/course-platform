package sqlite_database

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database/internal"
)

func (db *SQLiteDatabase) UserLikedTutorial(userId, slug string) (bool, error) {
	liked, err := internal.UserLikedTutorial(db.connection, userId, slug)
	if err != nil {
		db.ErrorLog.Printf("Failed to see if user liked tutorial %s in the database: %s\n", slug, err)
		return false, err
	}

	return liked, nil
}

func (db *SQLiteDatabase) UserLikeTutorial(userId, slug string) error {
	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate ID for new tutorials_likes database row: %s\n", err)
		return err
	}

	if err := internal.UserLikeTutorial(db.connection, id, userId, slug); err != nil {
		db.ErrorLog.Printf("Failed to insert new row into tutorials_likes database table: %s\n", err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) UserDislikeTutorial(userId, slug string) error {
	if err := internal.UserDislikeTutorial(db.connection, userId, slug); err != nil {
		db.ErrorLog.Printf("Failed to delete row from tutorials_likes database table: %s\n", err)
		return err
	}

	return nil
}
