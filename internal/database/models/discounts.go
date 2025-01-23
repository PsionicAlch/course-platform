package models

import "time"

// DiscountModel is a struct representation of the discounts table.
type DiscountModel struct {
	ID          string
	Title       string
	Description string
	Code        string
	Discount    uint
	Uses        uint
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
