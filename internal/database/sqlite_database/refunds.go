package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func (db *SQLiteDatabase) AdminGetRefunds(term string, status string, page, elements uint) ([]*models.RefundModel, error) {
	query := `SELECT r.id, r.user_id, r.course_purchase_id, r.refund_status, r.created_at, r.updated_at FROM refunds AS r LEFT JOIN users AS u ON r.user_id = u.id LEFT JOIN course_purchases AS cp ON r.course_purchase_id = cp.id LEFT JOIN courses AS c ON cp.course_id = c.id WHERE 1=1`
	args := []any{}

	if term != "" {
		query += " AND (LOWER(r.id) LIKE '%' || ? || '%' OR LOWER(u.name) LIKE '%' || ? || '%' OR LOWER(u.surname) LIKE '%' || ? || '%' OR LOWER(c.title) LIKE '%' || ? || '%')"
		args = append(args, term, term, term, term)
	}

	if status != "" {
		query += " AND r.refund_status = ?"
		args = append(args, status)
	}

	offset := (page - 1) * elements
	query += " ORDER BY r.updated_at DESC, r.created_at DESC LIMIT ? OFFSET ?;"
	args = append(args, elements, offset)

	var refunds []*models.RefundModel

	rows, err := db.connection.Query(query, args...)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all refunds from the database: %s\n", err)
		return nil, err
	}

	for rows.Next() {
		var refund models.RefundModel

		if err := rows.Scan(&refund.ID, &refund.UserID, &refund.CoursePurchaseID, &refund.RefundStatus, &refund.CreatedAt, &refund.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from refunds table: %s\n", err)
			return nil, err
		}

		refunds = append(refunds, &refund)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to get all refunds from the database: %s\n", err)
		return nil, err
	}

	return refunds, nil
}

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

	_, err := db.connection.Exec(query, status.String(), refundId)
	if err != nil {
		db.ErrorLog.Printf("Failed to update refund (\"%s\") status: %s\n", refundId, err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) CountRefunds() (uint, error) {
	query := `SELECT COUNT(id) FROM refunds;`

	var count uint

	row := db.connection.QueryRow(query)
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		db.ErrorLog.Printf("Failed to count the number of refunds in the database: %s\n", err)
		return 0, err
	}

	return count, nil
}
