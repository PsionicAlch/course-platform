package models

import "time"

// RefundModel is a struct representation of the refunds table.
type RefundModel struct {
	ID               string
	UserID           string
	CoursePurchaseID string
	RefundStatus     string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
