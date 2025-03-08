package sqlite_database

import (
	"github.com/PsionicAlch/course-platform/internal/database"
	"github.com/PsionicAlch/course-platform/internal/database/models"
	"github.com/PsionicAlch/course-platform/internal/database/sqlite_database/internal"
)

// RegisterAffiliatePointsChange adds a new row to the affiliate_points_history table.
func (db *SQLiteDatabase) RegisterAffiliatePointsChange(userId, courseId string, pointsChange int, reason string) error {
	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate ID for affiliate points history: %s\n", err)
		return err
	}

	if err := internal.RegisterAffiliatePointsChange(db.connection, id, userId, courseId, pointsChange, reason); err != nil {
		db.ErrorLog.Printf("Failed to save affiliate point change to database: %s\n", err)
		return err
	}

	return nil
}

// CountUserAffiliateHistory counts all the times the given user's affiliate code was used.
func (db *SQLiteDatabase) CountUserAffiliateHistory(userId string) (uint, error) {
	query := `SELECT COUNT(id) FROM affiliate_points_history WHERE user_id = ?;`

	var count uint

	row := db.connection.QueryRow(query, userId)
	if err := row.Scan(&count); err != nil {
		db.ErrorLog.Printf("Failed to count the number of affiliate_points_history rows connected to user (\"%s\"): %s\n", userId, err)
		return 0, err
	}

	return count, nil
}

// GetUserAffiliatePointsHistory gets a slice of AffiliatePointsHistoryModel for a given user.
func (db *SQLiteDatabase) GetUserAffiliatePointsHistory(userId string, page, elements uint) ([]*models.AffiliatePointsHistoryModel, error) {
	query := `SELECT id, user_id, course_id, points_change, reason, created_at FROM affiliate_points_history WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;;`

	var history []*models.AffiliatePointsHistoryModel

	offset := (page - 1) * elements

	rows, err := db.connection.Query(query, userId, elements, offset)
	if err != nil {
		db.ErrorLog.Printf("Failed to get user's (\"%s\") affiliate points history: %s\n", userId, err)
		return nil, err
	}

	for rows.Next() {
		var h models.AffiliatePointsHistoryModel

		if err := rows.Scan(&h.ID, &h.UserID, &h.CourseID, &h.PointsChange, &h.Reason, &h.CreatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from affiliate_points_history table: %s\n", err)
			return nil, err
		}

		history = append(history, &h)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to get user's (\"%s\") affiliate points history: %s\n", userId, err)
		return nil, err
	}

	return history, nil
}
