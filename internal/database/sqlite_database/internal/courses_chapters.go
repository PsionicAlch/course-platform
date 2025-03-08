package internal

import (
	"github.com/PsionicAlch/course-platform/internal/database"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// AddChapter adds a new chapter row to the database. This function works with either a database connection or a
// database transaction. This function will NOT throw an error upon a unique constraint violation.
func AddChapter(dbFacade SqlDbFacade, id, title, slug string, chapter int, content, fileChecksum, fileKey, courseKey string) error {
	query := `INSERT INTO course_chapters (id, title, slug, chapter, content, course_id, file_checksum, file_key) VALUES (?, ?, ?, ?, ?, (SELECT id FROM courses WHERE file_key = ?), ?, ?);`

	result, err := dbFacade.Exec(query, id, title, slug, chapter, content, courseKey, fileChecksum, fileKey)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return database.ErrChapterAlreadyExists
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNoRowsAffected
	}

	return nil
}

// UpdateChapter updates a chapter in the database based off the provided ID. This function works with either a database
// connection or a database transaction.
func UpdateChapter(dbFacade SqlDbFacade, id, title, slug string, chapter int, content, fileChecksum, fileKey, courseKey string) error {
	query := `UPDATE course_chapters SET title = ?, slug = ?, chapter = ?, content = ?, course_id = (SELECT id FROM courses WHERE file_key = ?), file_checksum = ?, file_key = ? WHERE id = ?;`

	result, err := dbFacade.Exec(query, title, slug, chapter, content, courseKey, fileChecksum, fileKey, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNoRowsAffected
	}

	return nil
}
