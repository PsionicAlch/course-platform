package internal

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

// AddCourse adds a new course row to the database. This function works with either a database connection or a database
// transaction.
func AddCourse(dbFacade SqlDbFacade, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey string) error {
	query := `INSERT INTO courses (id, title, slug, description, thumbnail_url, banner_url, content, file_checksum, file_key) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	results, err := dbFacade.Exec(query, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNoRowsAffected
	}

	return nil
}

// UpdateCourse updates the course row based on the provided ID. This function works with either a database connection
// or a database transaction.
func UpdateCourse(dbFacade SqlDbFacade, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey string) error {
	query := `UPDATE courses SET title = ?, slug = ?, description = ?, thumbnail_url = ?, banner_url = ?, content = ?, published = 0, file_checksum = ?, file_key = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;`
	results, err := dbFacade.Exec(query, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey, id)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNoRowsAffected
	}

	return nil
}

// GetCourseByID retrieves a CourseModel based on the provided course ID. This function works with either a database
// connection or a database transaction.
func GetCourseByID(dbFacade SqlDbFacade, courseId string) (*models.CourseModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM courses WHERE id = ?;`

	var course models.CourseModel
	var published int

	row := dbFacade.QueryRow(query, courseId)
	if err := row.Scan(&course.ID, &course.Title, &course.Slug, &course.Description, &course.ThumbnailURL, &course.BannerURL, &course.Content, &published, &course.AuthorID, &course.FileChecksum, &course.FileKey, &course.CreatedAt, &course.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	course.Published = published == 1

	return &course, nil
}
