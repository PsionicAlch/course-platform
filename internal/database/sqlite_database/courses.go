package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

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

func (db *SQLiteDatabase) GetCourses(term string, page, elements int) ([]*models.CourseModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM courses WHERE published = 1 AND author_id IS NOT NULL AND (LOWER(title) LIKE '%' || ? || '%' OR LOWER(description) LIKE '%' || ? || '%') ORDER BY updated_at DESC, title ASC LIMIT ? OFFSET ?;`

	offset := (page - 1) * elements

	var courses []*models.CourseModel

	rows, err := db.connection.Query(query, term, term, elements, offset)
	if err != nil {
		db.ErrorLog.Printf("Failed to query database for all courses on page %d: %s\n", page, err)
		return nil, err
	}

	for rows.Next() {
		var course models.CourseModel
		var published int

		if err := rows.Scan(&course.ID, &course.Title, &course.Slug, &course.Description, &course.ThumbnailURL, &course.BannerURL, &course.Content, &published, &course.AuthorID, &course.FileChecksum, &course.FileKey, &course.CreatedAt, &course.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from courses table on page %d: %s\n", page, err)
			return nil, err
		}

		course.Published = published == 1

		courses = append(courses, &course)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to query database for all courses on page %d: %s\n", page, err)
		return nil, err
	}

	return courses, nil
}

func (db *SQLiteDatabase) GetCourseByFileKey(fileKey string) (*models.CourseModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM courses WHERE file_key = ?;`

	var course models.CourseModel
	var published int

	row := db.connection.QueryRow(query, fileKey)
	if err := row.Scan(&course.ID, &course.Title, &course.Slug, &course.Description, &course.ThumbnailURL, &course.BannerURL, &course.Content, &published, &course.AuthorID, &course.FileChecksum, &course.FileKey, &course.CreatedAt, &course.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to get course by file key from the database: %s\n", err)
		return nil, err
	}

	course.Published = published == 1

	return &course, nil
}

func (db *SQLiteDatabase) GetCourseBySlug(slug string) (*models.CourseModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM courses WHERE slug = ?;`

	var course models.CourseModel
	var published int

	row := db.connection.QueryRow(query, slug)
	if err := row.Scan(&course.ID, &course.Title, &course.Slug, &course.Description, &course.ThumbnailURL, &course.BannerURL, &course.Content, &published, &course.AuthorID, &course.FileChecksum, &course.FileKey, &course.CreatedAt, &course.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to get course by slug from the database: %s\n", err)
		return nil, err
	}

	course.Published = published == 1

	return &course, nil
}
