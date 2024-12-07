package models

import (
	"database/sql"
	"time"
)

type CoursePurchaseModel struct {
	ID                      string
	UserID                  string
	CourseID                string
	PaymentKey              string
	StripeCheckoutSessionID string
	AffiliateCode           sql.NullString
	DiscountCode            sql.NullString
	AffiliatePointsUsed     uint
	AmountPaid              float64
	PaymentStatus           string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}
