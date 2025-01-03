package internal

import (
	"database/sql"
	"fmt"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
)

func HasUserPurchasedCourse(dbFacade SqlDbFacade, userId, courseId string) (bool, error) {
	query := `SELECT id FROM course_purchases WHERE user_id = ? AND course_id = ? AND payment_status = ?;`

	var id string

	row := dbFacade.QueryRow(query, userId, courseId, database.Succeeded.String())
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("User ID: %s\nCourse ID: %s\nPayment Status: %s\n", userId, courseId, database.Succeeded.String())
			return false, nil
		}

		return false, err
	}

	return id != "", nil
}

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
