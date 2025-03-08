package internal

import "github.com/PsionicAlch/course-platform/internal/database"

// RegisterAffiliatePointsChange adds a new row to the affiliate_points_history table in a way that works with a normal
// database connection or a database transaction.
func RegisterAffiliatePointsChange(dbFacade SqlDbFacade, id, userId, courseId string, pointsChange int, reason string) error {
	query := `INSERT INTO affiliate_points_history (id, user_id, course_id, points_change, reason) VALUES (?, ?, ?, ?, ?);`

	result, err := dbFacade.Exec(query, id, userId, courseId, pointsChange, reason)
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
