package models

import "time"

// AffiliatePointsHistoryModel is a struct representation of the affiliate_points_history table.
type AffiliatePointsHistoryModel struct {
	ID           string
	UserID       string
	CourseID     string
	PointsChange int
	Reason       string
	CreatedAt    time.Time
}
