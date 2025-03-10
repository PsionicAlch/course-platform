package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/course-platform/internal/database/models"
)

func (db *SQLiteDatabase) GetAllChapters() ([]*models.ChapterModel, error) {
	query := `SELECT id, title, slug, chapter, content, course_id, file_checksum, file_key, created_at, updated_at FROM course_chapters ORDER BY chapter ASC;`

	var chapters []*models.ChapterModel

	rows, err := db.connection.Query(query)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all chapters from the database: %s\n", err)
		return nil, err
	}

	for rows.Next() {
		var chapter models.ChapterModel

		if err := rows.Scan(&chapter.ID, &chapter.Title, &chapter.Slug, &chapter.Chapter, &chapter.Content, &chapter.CourseID, &chapter.FileChecksum, &chapter.FileKey, &chapter.CreatedAt, &chapter.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from chapters table: %s\n", err)
			return nil, err
		}

		chapters = append(chapters, &chapter)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to get all chapters from the database: %s\n", err)
		return nil, err
	}

	return chapters, nil
}

func (db *SQLiteDatabase) GetChapterBySlug(chapterSlug string) (*models.ChapterModel, error) {
	query := `SELECT id, title, slug, chapter, content, course_id, file_checksum, file_key, created_at, updated_at FROM course_chapters WHERE slug = ?;`

	var chapter models.ChapterModel

	row := db.connection.QueryRow(query, chapterSlug)
	if err := row.Scan(&chapter.ID, &chapter.Title, &chapter.Slug, &chapter.Chapter, &chapter.Content, &chapter.CourseID, &chapter.FileChecksum, &chapter.FileKey, &chapter.CreatedAt, &chapter.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to get chapter using slug (\"%s\") from the database: %s\n", chapterSlug, err)
		return nil, err
	}

	return &chapter, nil
}

func (db *SQLiteDatabase) GetChapterByFileKey(fileKey string) (*models.ChapterModel, error) {
	query := `SELECT id, title, slug, chapter, content, course_id, file_checksum, file_key, created_at, updated_at FROM course_chapters WHERE file_key = ?;`

	var chapter models.ChapterModel

	row := db.connection.QueryRow(query, fileKey)
	if err := row.Scan(&chapter.ID, &chapter.Title, &chapter.Slug, &chapter.Chapter, &chapter.Content, &chapter.CourseID, &chapter.FileChecksum, &chapter.FileKey, &chapter.CreatedAt, &chapter.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to get chapter using file key from the database: %s\n", err)
		return nil, err
	}

	return &chapter, nil
}

func (db *SQLiteDatabase) CountChapters(courseId string) (int, error) {
	query := `SELECT COUNT(id) FROM course_chapters WHERE course_id = ?;`

	var chapters int

	row := db.connection.QueryRow(query, courseId)
	if err := row.Scan(&chapters); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		db.ErrorLog.Printf("Failed to count all the chapters for course with ID \"%s\": %s\n", courseId, err)
		return -1, err
	}

	return chapters, nil
}

func (db *SQLiteDatabase) GetCourseChapters(courseId string) ([]*models.ChapterModel, error) {
	query := `SELECT id, title, slug, chapter, content, course_id, file_checksum, file_key, created_at, updated_at FROM course_chapters WHERE course_id = ? ORDER BY chapter ASC;`

	var chapters []*models.ChapterModel

	rows, err := db.connection.Query(query, courseId)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all chapters for course (\"%s\") from the database: %s\n", courseId, err)
		return nil, err
	}

	for rows.Next() {
		var chapter models.ChapterModel

		if err := rows.Scan(&chapter.ID, &chapter.Title, &chapter.Slug, &chapter.Chapter, &chapter.Content, &chapter.CourseID, &chapter.FileChecksum, &chapter.FileKey, &chapter.CreatedAt, &chapter.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from chapters table: %s\n", err)
			return nil, err
		}

		chapters = append(chapters, &chapter)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to get all chapters for course (\"%s\") from the database: %s\n", courseId, err)
		return nil, err
	}

	return chapters, nil
}
