package sqlite_database

import "github.com/PsionicAlch/psionicalch-home/internal/database/models"

func (db *SQLiteDatabase) GetAllCourses() ([]*models.CourseModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM courses;`

	var courses []*models.CourseModel

	rows, err := db.connection.Query(query)
	if err != nil {
		db.ErrorLog.Printf("Failed to query database for all courses: %s\n", err)
		return nil, err
	}

	for rows.Next() {
		var course models.CourseModel
		var published int

		if err := rows.Scan(&course.ID, &course.Title, &course.Slug, &course.Description, &course.ThumbnailURL, &course.BannerURL, &course.Content, &published, &course.AuthorID, &course.FileChecksum, &course.FileKey, &course.CreatedAt, &course.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from courses table: %s\n", err)
			return nil, err
		}

		course.Published = published == 1

		courses = append(courses, &course)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to query database for all courses: %s\n", err)
		return nil, err
	}

	return courses, nil
}
