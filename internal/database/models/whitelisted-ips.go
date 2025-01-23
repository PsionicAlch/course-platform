package models

import "time"

// WhitelistedIPModels is a struct representation of the whitelisted_ips table.
type WhitelistedIPModel struct {
	ID        string
	UserID    string
	IPAddress string
	CreatedAt time.Time
}
