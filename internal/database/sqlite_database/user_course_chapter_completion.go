package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/course-platform/internal/database"
	"github.com/PsionicAlch/course-platform/internal/database/models"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func (db *SQLiteDatabase) HasUserCompletedChapter(userId, courseId, chapterId string) (bool, error) {
	query := `SELECT id FROM user_course_chapter_completion WHERE user_id = ? AND course_id = ? AND chapter_id = ?;`

	var id string

	row := db.connection.QueryRow(query, userId, courseId, chapterId)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		db.ErrorLog.Printf("Failed to check if user (\"%s\") has finished course (\"%s\") chapter (\"%s\"): %s\n", userId, courseId, chapterId, err)
		return false, err
	}

	return id != "", nil
}

func (db *SQLiteDatabase) GetAllChaptersCompleted(userId, courseId string) ([]*models.ChapterModel, error) {
	query := `SELECT cc.id, cc.title, cc.slug, cc.chapter, cc.content, cc.course_id, cc.file_checksum, cc.file_key, cc.created_at, cc.updated_at FROM user_course_chapter_completion AS uccc LEFT JOIN course_chapters AS cc ON uccc.chapter_id = cc.id WHERE uccc.user_id = ? AND uccc.course_id = ? ORDER BY cc.chapter ASC;`

	var chapters []*models.ChapterModel

	rows, err := db.connection.Query(query, userId, courseId)
	if err != nil {
		db.ErrorLog.Printf("Failed to get course (\"%s\") chapters completed by user (\"%s\"): %s\n", courseId, userId, err)
		return nil, err
	}

	for rows.Next() {
		var chapter models.ChapterModel

		if err := rows.Scan(&chapter.ID, &chapter.Title, &chapter.Slug, &chapter.Chapter, &chapter.Content, &chapter.CourseID, &chapter.FileChecksum, &chapter.FileKey, &chapter.CreatedAt, &chapter.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read chapter from database: %s\n", err)
			return nil, err
		}

		chapters = append(chapters, &chapter)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to get course (\"%s\") chapters completed by user (\"%s\"): %s\n", courseId, userId, err)
		return nil, err
	}

	return chapters, nil
}

func (db *SQLiteDatabase) GetAllChaptersNotCompleted(userId, courseId string) ([]*models.ChapterModel, error) {
	query := `SELECT id, title, slug, chapter, content, course_id, file_checksum, file_key, created_at, updated_at FROM course_chapters WHERE course_id = ? EXCEPT SELECT cc.id, cc.title, cc.slug, cc.chapter, cc.content, cc.course_id, cc.file_checksum, cc.file_key, cc.created_at, cc.updated_at FROM course_chapters AS cc LEFT JOIN user_course_chapter_completion AS uccc ON cc.id = uccc.chapter_id WHERE uccc.user_id = ? AND cc.course_id = ? ORDER BY chapter ASC;`

	var chapters []*models.ChapterModel

	rows, err := db.connection.Query(query, courseId, userId, courseId)
	if err != nil {
		db.ErrorLog.Printf("Failed to get course (\"%s\") chapters not completed by user (\"%s\"): %s\n", courseId, userId, err)
		return nil, err
	}

	for rows.Next() {
		var chapter models.ChapterModel

		if err := rows.Scan(&chapter.ID, &chapter.Title, &chapter.Slug, &chapter.Chapter, &chapter.Content, &chapter.CourseID, &chapter.FileChecksum, &chapter.FileKey, &chapter.CreatedAt, &chapter.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read chapter from database: %s\n", err)
			return nil, err
		}

		chapters = append(chapters, &chapter)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to get course (\"%s\") chapters not completed by user (\"%s\"): %s\n", courseId, userId, err)
		return nil, err
	}

	return chapters, nil
}

func (db *SQLiteDatabase) FinishChapter(userId, chapterId, courseId string) error {
	query := `INSERT INTO user_course_chapter_completion (id, user_id, course_id, chapter_id) VALUES (?, ?, ?, ?);`

	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate ID for new user_course_chapter_completion table row: %s\n", err)
		return err
	}

	result, err := db.connection.Exec(query, id, userId, courseId, chapterId)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return nil
		}

		db.ErrorLog.Printf("Failed to insert new row to user_course_chapter_completion table: %s\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to get rows affected after inserting new row to user_course_chapter_completion table: %s\n", err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Println("No rows were affected after inserting new row to user_course_chapter_completion table")
		return database.ErrNoRowsAffected
	}

	return nil
}
