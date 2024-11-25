package internal

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func AddChapter(dbFacade SqlDbFacade, id, title string, chapter int, content, fileChecksum, fileKey, courseKey string) error {
	query := `INSERT INTO courses_chapters (id, title, chapter, content, course_id, file_checksum, file_key) VALUES (?, ?, ?, ?, (SELECT id FROM courses WHERE file_key = ?), ?, ?);`

	result, err := dbFacade.Exec(query, id, title, chapter, content, courseKey, fileChecksum, fileKey)
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

func UpdateChapter(dbFacade SqlDbFacade, id, title string, chapter int, content, fileChecksum, fileKey, courseKey string) error {
	query := `UPDATE courses_chapters SET title = ?, chapter = ?, content = ?, course_id = (SELECT id FROM courses WHERE file_key = ?), file_checksum = ?, file_key = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;`

	result, err := dbFacade.Exec(query, title, chapter, content, courseKey, fileChecksum, fileKey, id)
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
