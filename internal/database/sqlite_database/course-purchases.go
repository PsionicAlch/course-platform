package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
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

func (db *SQLiteDatabase) RegisterCoursePurchase(userId, courseId, stripeCheckoutSessionId string, affiliateCode, discountCode sql.NullString, affiliatePointsUsed int64, amountPaid float64) error {
	purchaseId, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate ID for new course purchase: %s\n", err)
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

	purchased, err := internal.HasUserPurchasedCourse(tx, user.ID, courseId)
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

	if user.AffiliatePoints < uint(affiliatePointsUsed) {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after error occurred: %s\n", err)
		}

		return database.ErrInsufficientAffiliatePoints
	}

	if err := internal.UpdateAffiliatePoints(tx, user.ID, user.AffiliatePoints-uint(affiliatePointsUsed)); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to update user's (\"%s\") affiliate points: %s\n", user.ID, err)
		return err
	}

	if err := internal.AddNewCoursePurchase(tx, purchaseId, user.ID, courseId, stripeCheckoutSessionId, affiliateCode, discountCode, affiliatePointsUsed, amountPaid); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to register course purchase: %s\n", err)
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
