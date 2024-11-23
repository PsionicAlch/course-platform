package sqlite_database

import "github.com/PsionicAlch/psionicalch-home/internal/database/models"

func (db *SQLiteDatabase) GetAllChapters() ([]*models.ChapterModel, error) {
	query := `SELECT id, title, chapter, content, course_id, file_checksum, file_key, created_at, updated_at FROM courses_chapters;`

	var chapters []*models.ChapterModel

	rows, err := db.connection.Query(query)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all chapters from the database: %s\n", err)
		return nil, err
	}

	for rows.Next() {
		var chapter models.ChapterModel

		if err := rows.Scan(&chapter.ID, &chapter.Title, &chapter.Chapter, &chapter.Content, &chapter.CourseID, &chapter.FileChecksum, &chapter.FileKey, &chapter.CreatedAt, &chapter.UpdatedAt); err != nil {
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
