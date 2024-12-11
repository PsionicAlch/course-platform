package models

import "time"

type AffiliatePointsHistoryModel struct {
	ID           string
	UserID       string
	CourseID     string
	PointsChange int
	Reason       string
	CreatedAt    time.Time
}
