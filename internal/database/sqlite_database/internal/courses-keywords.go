package internal

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func AddKeywordToCourse(dbFacade SqlDbFacade, id, keywordId, courseId string) error {
	query := `INSERT INTO courses_keywords (id, course_id, keyword_id) VALUES (?, ?, ?);`

	result, err := dbFacade.Exec(query, id, courseId, keywordId)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return nil
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

func DeleteAllKeywordsFromCourses(dbFacade SqlDbFacade, courseId string) error {
	query := `DELETE FROM courses_keywords WHERE course_id = ?;`

	_, err := dbFacade.Exec(query, courseId)
	if err != nil {
		return err
	}

	return nil
}
