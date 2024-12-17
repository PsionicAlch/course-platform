package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database/internal"
)

func (db *SQLiteDatabase) AdminGetCourses(term string, published *bool, authorId *string, boughtBy, keyword string, page, elements uint) ([]*models.CourseModel, error) {
	query := `SELECT DISTINCT c.id, c.title, c.slug, c.description, c.thumbnail_url, c.banner_url, c.content, c.published, c.author_id, c.file_checksum, c.file_key, c.created_at, c.updated_at FROM courses AS c LEFT JOIN course_purchases AS cp ON cp.course_id = c.id LEFT JOIN courses_keywords AS ck ON ck.course_id = c.id LEFT JOIN keywords AS k ON k.id = ck.keyword_id WHERE (LOWER(c.id) LIKE '%' || ? || '%' OR LOWER(c.title) LIKE '%' || ? || '%' OR LOWER(c.slug) LIKE '%' || ? || '%' OR LOWER(c.description) LIKE '%' || ? || '%' OR LOWER(k.keyword) LIKE '%' || ? || '%')`

	args := []any{term, term, term, term, term}

	if published != nil {
		query += " AND c.published = ?"

		var pubInt int
		if *published {
			pubInt = 1
		} else {
			pubInt = 0
		}

		args = append(args, pubInt)
	}

	if authorId != nil {
		if *authorId != "" {
			query += " AND c.author_id = ?"
			args = append(args, *authorId)
		}
	} else {
		query += " AND c.author_id IS NULL"
	}

	if boughtBy != "" {
		query += " AND cp.user_id = ?"
		args = append(args, boughtBy)
	}

	if keyword != "" {
		query += " AND k.keyword LIKE '%' || ? || '%'"
		args = append(args, keyword)
	}

	offset := (page - 1) * elements
	query += " ORDER BY c.created_at DESC, c.updated_at DESC, c.title ASC LIMIT ? OFFSET ?;"
	args = append(args, elements, offset)

	var courses []*models.CourseModel

	rows, err := db.connection.Query(query, args...)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all courses from the database: %s\n", err)
		db.ErrorLog.Printf("\nSQL Query Used: \n%s\n", query)

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
		db.ErrorLog.Printf("Failed to get all courses from the database: %s\n", err)
		db.ErrorLog.Printf("\nSQL Query Used: \n%s\n", query)

		return nil, err
	}

	return courses, nil
}

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

func (db *SQLiteDatabase) GetCourses(term string, authorId string, page, elements int) ([]*models.CourseModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM courses WHERE published = 1 AND author_id IS NOT NULL`
	args := []any{}

	if term != "" {
		query += " AND (LOWER(title) LIKE '%' || ? || '%' OR LOWER(description) LIKE '%' || ? || '%')"
		args = append(args, term, term)
	}

	if authorId != "" {
		query += " AND author_id = ?"
		args = append(args, authorId)
	}

	offset := (page - 1) * elements
	query += " ORDER BY updated_at DESC, title ASC LIMIT ? OFFSET ?;"
	args = append(args, elements, offset)

	var courses []*models.CourseModel

	rows, err := db.connection.Query(query, args...)
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

		db.ErrorLog.Printf("Failed to get course by slug (\"%s\") from the database: %s\n", slug, err)
		return nil, err
	}

	course.Published = published == 1

	return &course, nil
}

func (db *SQLiteDatabase) GetCourseByID(courseId string) (*models.CourseModel, error) {
	course, err := internal.GetCourseByID(db.connection, courseId)
	if err != nil {
		db.ErrorLog.Printf("Failed to get course by ID (\"%s\"): %s\n", courseId, err)
		return nil, err
	}

	return course, nil
}

func (db *SQLiteDatabase) CountCourses() (uint, error) {
	query := `SELECT COUNT(id) FROM courses;`

	var count uint

	row := db.connection.QueryRow(query)
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		db.ErrorLog.Printf("Failed to count all the courses in the database: %s\n", err)
		return 0, err
	}

	return count, nil
}

func (db *SQLiteDatabase) CountCoursesWrittenBy(authorId string) (uint, error) {
	query := `SELECT COUNT(id) FROM courses WHERE author_id = ?;`

	var count uint

	row := db.connection.QueryRow(query, authorId)
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		db.ErrorLog.Printf("Failed to count all the courses written by \"%s\": %s\n", authorId, err)
		return 0, err
	}

	return count, nil
}

func (db *SQLiteDatabase) PublishCourse(courseId string) error {
	query := `UPDATE courses SET published = 1 WHERE id = ?;`

	_, err := db.connection.Exec(query, courseId)
	if err != nil {
		db.ErrorLog.Printf("Failed to publish course \"%s\": %s\n", courseId, err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) UnpublishCourse(courseId string) error {
	query := `UPDATE courses SET published = 0 WHERE id = ?;`

	_, err := db.connection.Exec(query, courseId)
	if err != nil {
		db.ErrorLog.Printf("Failed to unpublish course \"%s\": %s\n", courseId, err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) UpdateCourseAuthor(courseId, authorId string) error {
	query := `UPDATE courses SET author_id = ? WHERE id = ?;`

	if authorId == "" {
		authorId = "NULL"
	}

	_, err := db.connection.Exec(query, authorId, courseId)
	if err != nil {
		db.ErrorLog.Printf("Failed to update course's (\"%s\") author (\"%s\"): %s\n", courseId, authorId, err)
		return err
	}

	return nil
}
