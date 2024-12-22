package sqlite_database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database/internal"
)

func (db *SQLiteDatabase) HasUserPurchasedCourse(userId, courseId string) (bool, error) {
	b, err := internal.HasUserPurchasedCourse(db.connection, userId, courseId)
	if err != nil {
		db.ErrorLog.Printf("Failed to check if user (\"%s\") has purchased course (\"%s\"): %s\n", userId, courseId, err)

		return false, err
	}

	return b, nil
}

func (db *SQLiteDatabase) RegisterCoursePurchase(userId, courseId, paymentKey, stripeCheckoutSessionId string, affiliateCode, discountCode sql.NullString, affiliatePointsUsed uint, amountPaid float64, token, tokenType string, validUntil time.Time) error {
	purchaseId, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate ID for new course purchase: %s\n", err)
		return err
	}

	affiliateHistoryId, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate ID for new affiliate points history: %s\n", err)
		return err
	}

	paymentTokenId, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate ID for new payment token: %s\n", err)
		return err
	}

	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("Failed to start new database transaction: %s\n", err)
		return err
	}

	user, err := internal.GetUserByID(tx, userId, database.All)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to get user by ID (\"%s\"): %s\n", userId, err)
		return err
	}

	course, err := internal.GetCourseByID(tx, courseId)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to get course by ID (\"%s\"): %s\n", courseId, err)
	}

	purchased, err := internal.HasUserPurchasedCourse(tx, user.ID, course.ID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to check if user (\"%s\") has already purchased this course (\"%s\"): %s\n", userId, courseId, err)
		return err
	}

	if purchased {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after error occurred: %s\n", err)
		}

		return database.ErrCourseAlreadyOwned
	}

	if user.AffiliatePoints < int(affiliatePointsUsed) {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after error occurred: %s\n", err)
		}

		return database.ErrInsufficientAffiliatePoints
	}

	if affiliatePointsUsed > 0 {
		if err := internal.RegisterAffiliatePointsChange(tx, affiliateHistoryId, user.ID, courseId, -1*int(affiliatePointsUsed), fmt.Sprintf("Purchased \"%s\"", course.Title)); err != nil {
			if err := tx.Rollback(); err != nil {
				db.ErrorLog.Printf("Failed to rollback changes after error occurred: %s\n", err)
			}

			db.ErrorLog.Printf("Failed to register affiliate point change: %s\n", err)
			return err
		}
	}

	if err := internal.AddNewCoursePurchase(tx, purchaseId, user.ID, courseId, paymentKey, stripeCheckoutSessionId, affiliateCode, discountCode, affiliatePointsUsed, amountPaid); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to register course purchase: %s\n", err)
		return err
	}

	if err := internal.AddToken(tx, paymentTokenId, token, tokenType, user.ID, validUntil); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to save the payment token: %s\n", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to commit changes after registering course purchase: %s\n", err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) CountCoursesWhereDiscountWasUsed(discountCode string) (uint, error) {
	query := `SELECT COUNT(id) FROM course_purchases WHERE discount_code = ?;`

	var count uint

	row := db.connection.QueryRow(query, discountCode)
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		db.ErrorLog.Printf("Failed to count the number of courses bought using discount code (\"%s\"): %s\n", discountCode, err)
		return 0, err
	}

	return count, nil
}

func (db *SQLiteDatabase) CountUsersWhoBoughtCourse(courseId string) (uint, error) {
	query := `SELECT COUNT(id) FROM course_purchases WHERE course_id = ?;`

	var count uint

	row := db.connection.QueryRow(query, courseId)
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		db.ErrorLog.Printf("Failed to count the number times course \"%s\" was purchased: %s\n", courseId, err)
		return 0, err
	}

	return count, nil
}

func (db *SQLiteDatabase) GetCoursePurchaseByPaymentKey(paymentKey string) (*models.CoursePurchaseModel, error) {
	query := `SELECT id, user_id, course_id, payment_key, stripe_checkout_session_id, affiliate_code, discount_code, affiliate_points_used, amount_paid, payment_status, created_at, updated_at FROM course_purchases WHERE payment_key = ?;`

	var coursePurchase models.CoursePurchaseModel

	row := db.connection.QueryRow(query, paymentKey)
	if err := row.Scan(&coursePurchase.ID, &coursePurchase.UserID, &coursePurchase.CourseID, &coursePurchase.PaymentKey, &coursePurchase.StripeCheckoutSessionID, &coursePurchase.AffiliateCode, &coursePurchase.DiscountCode, &coursePurchase.AffiliatePointsUsed, &coursePurchase.AmountPaid, &coursePurchase.PaymentStatus, &coursePurchase.CreatedAt, &coursePurchase.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to find course purchase by payment key (\"%s\"): %s\n", paymentKey, err)
		return nil, err
	}

	return &coursePurchase, nil
}

func (db *SQLiteDatabase) UpdateCoursePurchasePaymentStatus(coursePurchaseId string, status database.PaymentStatus) error {
	query := `UPDATE course_purchases SET payment_status = ? WHERE id = ?;`

	result, err := db.connection.Exec(query, status.String(), coursePurchaseId)
	if err != nil {
		db.ErrorLog.Printf("Failed to update course purchase's (\"%s\") payment status: %s\n", coursePurchaseId, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to query database for rows affected after updating course purchase's (\"%s\") payment status: %s\n", coursePurchaseId, err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Printf("No rows were affected after updating course purchase's (\"%s\") payment status\n", coursePurchaseId)
		return database.ErrNoRowsAffected
	}

	return nil
}

func (db *SQLiteDatabase) GetCoursesBoughtByUser(term, userId string, page, elements uint) ([]*models.CourseModel, error) {
	query := `SELECT c.id, c.title, c.slug, c.description, c.thumbnail_url, c.banner_url, c.content, c.published, c.author_id, c.file_checksum, c.file_key, c.created_at, c.updated_at FROM course_purchases AS cp JOIN courses AS c ON cp.course_id = c.id WHERE cp.user_id = ?`
	args := []any{userId}

	if term != "" {
		query += " AND (LOWER(c.title) LIKE '%' || ? || '%' OR LOWER(c.slug) LIKE '%' || ? || '%' OR LOWER(c.description) LIKE '%' || ? || '%')"
		args = append(args, term, term, term)
	}

	offset := (page - 1) * elements
	query += " ORDER BY cp.updated_at DESC LIMIT ? OFFSET ?;"
	args = append(args, elements, offset)

	var courses []*models.CourseModel

	rows, err := db.connection.Query(query, args...)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all courses purchased bought by user (\"%s\"): %s\n", userId, err)
		return nil, err
	}

	for rows.Next() {
		var course models.CourseModel
		var published int

		if err := rows.Scan(&course.ID, &course.Title, &course.Slug, &course.Description, &course.ThumbnailURL, &course.BannerURL, &course.Content, &published, &course.AuthorID, &course.FileChecksum, &course.FileKey, &course.CreatedAt, &course.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read course from the database: %s\n", err)
			return nil, err
		}

		course.Published = published == 1

		courses = append(courses, &course)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to get all courses purchased bought by user (\"%s\"): %s\n", userId, err)
		return nil, err
	}

	return courses, nil
}

func (db *SQLiteDatabase) GetAllCoursesBoughtByUser(userId string) ([]*models.CourseModel, error) {
	query := `SELECT c.id, c.title, c.slug, c.description, c.thumbnail_url, c.banner_url, c.content, c.published, c.author_id, c.file_checksum, c.file_key, c.created_at, c.updated_at FROM course_purchases AS cp LEFT JOIN courses AS c ON cp.course_id = c.id WHERE cp.user_id = ? AND cp.payment_status = 'Succeeded' ORDER BY cp.updated_at DESC;`

	var courses []*models.CourseModel

	rows, err := db.connection.Query(query, userId)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all courses purchased bought by user (\"%s\"): %s\n", userId, err)
		return nil, err
	}

	for rows.Next() {
		var course models.CourseModel
		var published int

		if err := rows.Scan(&course.ID, &course.Title, &course.Slug, &course.Description, &course.ThumbnailURL, &course.BannerURL, &course.Content, &published, &course.AuthorID, &course.FileChecksum, &course.FileKey, &course.CreatedAt, &course.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read course from the database: %s\n", err)
			return nil, err
		}

		course.Published = published == 1

		courses = append(courses, &course)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to get all courses purchased bought by user (\"%s\"): %s\n", userId, err)
		return nil, err
	}

	return courses, nil
}
