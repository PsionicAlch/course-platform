package internal

import (
	"database/sql"

	"github.com/PsionicAlch/course-platform/internal/database"
)

// HasUserPurchasedCourse checks if there is a database row that indicates the user has purchased
// the provided course. This function works with normal database connections or database transactions.
func HasUserPurchasedCourse(dbFacade SqlDbFacade, userId, courseId string) (bool, error) {
	query := `SELECT id FROM course_purchases WHERE user_id = ? AND course_id = ? AND payment_status = ?;`

	var id string

	row := dbFacade.QueryRow(query, userId, courseId, database.Succeeded.String())
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return id != "", nil
}

// AddNewCoursePurchase adds a new course purchase row in the database. This function works with normal database
// connections or database transactions.
func AddNewCoursePurchase(dbFacade SqlDbFacade, purchaseId, userId, courseId, paymentKey, stripeCheckoutSessionId string, affiliateCode, discountCode sql.NullString, affiliatePointsUsed uint, amountPaid float64) error {
	query := `INSERT INTO course_purchases (id, user_id, course_id, payment_key, stripe_checkout_session_id, affiliate_code, discount_code, affiliate_points_used, amount_paid) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	result, err := dbFacade.Exec(query, purchaseId, userId, courseId, paymentKey, stripeCheckoutSessionId, affiliateCode, discountCode, affiliatePointsUsed, amountPaid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNoRowsAffected
	}

	return nil
}
