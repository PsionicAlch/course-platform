package models

import "time"

type RefundModel struct {
	ID               string
	UserID           string
	CoursePurchaseID string
	RefundStatus     string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
