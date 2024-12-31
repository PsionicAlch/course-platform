package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func (db *SQLiteDatabase) RegisterRefund(userId, coursePurchaseId string, status database.RefundStatus) error {
	query := `INSERT INTO refunds (id, user_id, course_purchase_id, refund_status) VALUES (?, ?, ?, ?);`

	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate ID for new refund: %s\n", err)
		return err
	}

	result, err := db.connection.Exec(query, id, userId, coursePurchaseId, status.String())
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return nil
		}

		db.ErrorLog.Printf("Failed to insert new refund: %s\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to get rows affected after inserting new refund: %s\n", err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Println("No rows affected after inserting new refund")
		return database.ErrNoRowsAffected
	}

	return nil
}

func (db *SQLiteDatabase) GetRefundWithCoursePurchaseID(coursePurchaseId string) (*models.RefundModel, error) {
	query := `SELECT id, user_id, course_purchase_id, refund_status, created_at, updated_at FROM refunds WHERE course_purchase_id = ?;`

	var refund models.RefundModel

	row := db.connection.QueryRow(query, coursePurchaseId)
	if err := row.Scan(&refund.ID, &refund.UserID, &refund.CoursePurchaseID, &refund.RefundStatus, &refund.CreatedAt, &refund.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to get refund from course purchase ID (\"%s\"): %s\n", coursePurchaseId, err)
		return nil, err
	}

	return &refund, nil
}

func (db *SQLiteDatabase) UpdateRefundStatus(refundId string, status database.RefundStatus) error {
	query := `UPDATE refunds SET refund_status = ? WHERE id = ?;`

	result, err := db.connection.Exec(query, status.String(), refundId)
	if err != nil {
		db.ErrorLog.Printf("Failed to update refund (\"%s\") status: %s\n", refundId, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to get rows affected after updating refund status: %s\n", err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Println("No rows were affected after updating refund status")
		return database.ErrNoRowsAffected
	}

	return nil
}
