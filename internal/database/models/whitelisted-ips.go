package models

import "time"

type WhitelistedIPModel struct {
	ID        string
	UserID    string
	IPAddress string
	CreatedAt time.Time
}
