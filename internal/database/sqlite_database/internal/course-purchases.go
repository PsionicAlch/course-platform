package internal

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
)

func HasUserPurchasedCourse(dbFacade SqlDbFacade, userId, courseId string) (bool, error) {
	query := `SELECT id FROM course_purchases WHERE user_id = ? AND course_id = ? AND payment_status = 'Succeeded';`

	var id string

	row := dbFacade.QueryRow(query, userId, courseId)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return id != "", nil
}

func AddNewCoursePurchase(dbFacade SqlDbFacade, purchaseId, userId, courseId, stripeCheckoutSessionId string, affiliateCode, discountCode sql.NullString, affiliatePointsUsed int64, amountPaid float64) error {
	query := `INSERT INTO course_purchases (id, user_id, course_id, stripe_checkout_session_id, affiliate_code, discount_code, affiliate_points_used, amount_paid) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

	result, err := dbFacade.Exec(query, purchaseId, userId, courseId, stripeCheckoutSessionId, affiliateCode, discountCode, affiliatePointsUsed, amountPaid)
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
