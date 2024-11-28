package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
)

func (db *SQLiteDatabase) UserBookmarkedTutorial(userId, slug string) (bool, error) {
	query := `SELECT t.id FROM tutorials_bookmarks AS tb JOIN tutorials AS t ON tb.tutorial_id = t.id WHERE tb.user_id = ? AND t.slug = ?;`

	var id string

	row := db.connection.QueryRow(query, userId, slug)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		db.ErrorLog.Printf("Failed to query row from tutorials_bookmarks table in the database: %s\n", err)
		return false, err
	}

	return id != "", nil
}

func (db *SQLiteDatabase) UserBookmarkTutorial(userId, slug string) error {
	query := `INSERT INTO tutorials_bookmarks (id, user_id, tutorial_id) VALUES (?, ?, (SELECT id FROM tutorials WHERE slug = ?));`

	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate ID for new tutorials_bookmarks database row: %s\n", err)
		return err
	}

	result, err := db.connection.Exec(query, id, userId, slug)
	if err != nil {
		db.ErrorLog.Printf("Failed to add new row to tutorials_bookmarks table in the database: %s\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to query for the number of rows affected after adding new row to tutorials_bookmarks table: %s\n", err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Println("No rows were affected after adding new row to tutorials_bookmarks table in the database.")
		return database.ErrNoRowsAffected
	}

	return nil
}

func (db *SQLiteDatabase) UserUnbookmarkTutorial(userId, slug string) error {
	query := `DELETE FROM tutorials_bookmarks WHERE user_id = ? AND tutorial_id = (SELECT id FROM tutorials WHERE slug = ?);`

	result, err := db.connection.Exec(query, userId, slug)
	if err != nil {
		db.ErrorLog.Printf("Failed to delete row from tutorials_bookmarks table in the database: %s\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to query the number of rows affected after deleting row from tutorial_bookmarks table in the database: %s\n", err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Println("No rows were affected after deleting a row from the tutorials_bookmarks table in the database")
		return database.ErrNoRowsAffected
	}

	return nil
}

func (db *SQLiteDatabase) CountTutorialBookmarks(tutorialId string) (uint, error) {
	query := `SELECT COUNT(id) FROM tutorials_bookmarks WHERE tutorial_id = ?;`

	var count uint

	row := db.connection.QueryRow(query, tutorialId)
	if err := row.Scan(&count); err != nil {
		db.ErrorLog.Printf("Failed to count the amount of times the tutorial \"%s\" has been bookmarked: %s\n", tutorialId, err)
		return 0, err
	}

	return count, nil
}
