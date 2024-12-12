package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
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
