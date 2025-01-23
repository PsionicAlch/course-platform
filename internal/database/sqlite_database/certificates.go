package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// AddCertificate adds a new certificate row the database.
func (db *SQLiteDatabase) AddCertificate(userId, courseId string) error {
	query := `INSERT INTO certificates (id, user_id, course_id) VALUES (?, ?, ?);`

	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate ID for new certificate: %s\n", err)
		return err
	}

	result, err := db.connection.Exec(query, id, userId, courseId)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return nil
		}

		db.ErrorLog.Printf("Failed to insert new certificate into the database: %s\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to get rows affected after inserting new certificate: %s\n", err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Println("No rows were affected after inserting new certificate")
		return database.ErrNoRowsAffected
	}

	return nil
}

// GetCertificateFromID retrieves a CertificateModel from the database by the given ID.
func (db *SQLiteDatabase) GetCertificateFromID(certificateId string) (*models.CertificateModel, error) {
	query := `SELECT id, user_id, course_id, created_at FROM certificates WHERE id = ?;`

	var certificate models.CertificateModel

	row := db.connection.QueryRow(query, certificateId)
	if err := row.Scan(&certificate.ID, &certificate.UserID, &certificate.CourseID, &certificate.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to get certificate (\"%s\") from the database: %s\n", certificateId, err)
		return nil, err
	}

	return &certificate, nil
}

// GetCertificateFromUserAndCourse retrieves a CertificateModel based on the user ID and course ID.
func (db *SQLiteDatabase) GetCertificateFromUserAndCourse(userId, courseId string) (*models.CertificateModel, error) {
	query := `SELECT id, user_id, course_id, created_at FROM certificates WHERE user_id = ? AND course_id = ?;`

	var certificate models.CertificateModel

	row := db.connection.QueryRow(query, userId, courseId)
	if err := row.Scan(&certificate.ID, &certificate.UserID, &certificate.CourseID, &certificate.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to get certificate from database for user (\"%s\") and course (\"%s\"): %s\n", userId, courseId, err)
		return nil, err
	}

	return &certificate, nil
}

// GetUserFromCertificate retrieves a UserModel from the given certificate ID.
func (db *SQLiteDatabase) GetUserFromCertificate(certificateId string) (*models.UserModel, error) {
	query := `SELECT u.id, u.name, u.surname, u.slug, u.email, u.password, u.is_admin, u.is_author, u.affiliate_code, u.affiliate_points, u.created_at, u.updated_at FROM certificates AS c LEFT JOIN users AS u ON c.user_id = u.id WHERE c.id = ?;`

	var user models.UserModel
	var isAdmin int
	var isAuthor int

	row := db.connection.QueryRow(query, certificateId)
	if err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Slug, &user.Email, &user.Password, &isAdmin, &isAuthor, &user.AffiliateCode, &user.AffiliatePoints, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to get user from certificate (\"%s\"): %s\n", certificateId, err)
		return nil, err
	}

	user.IsAdmin = isAdmin == 1
	user.IsAuthor = isAuthor == 1

	return &user, nil
}

// GetCourseFromCertificate retrieves a CourseModel from the given certificate ID.
func (db *SQLiteDatabase) GetCourseFromCertificate(certificateId string) (*models.CourseModel, error) {
	query := `SELECT c.id, c.title, c.slug, c.description, c.thumbnail_url, c.banner_url, c.content, c.published, c.author_id, c.file_checksum, c.file_key, c.created_at, c.updated_at FROM certificates AS cf LEFT JOIN courses AS c ON cf.course_id = c.id WHERE cf.id = ?;`

	var course models.CourseModel
	var published int

	row := db.connection.QueryRow(query, certificateId)
	if err := row.Scan(&course.ID, &course.Title, &course.Slug, &course.Description, &course.ThumbnailURL, &course.BannerURL, &course.Content, &published, &course.AuthorID, &course.FileChecksum, &course.FileKey, &course.CreatedAt, &course.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to get course from certificate (\"%s\"): %s\n", certificateId, err)
		return nil, err
	}

	course.Published = published == 1

	return &course, nil
}
