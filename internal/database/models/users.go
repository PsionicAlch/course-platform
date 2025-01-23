package models

import "time"

// UserModel is a struct representation of the users table.
type UserModel struct {
	ID              string
	Name            string
	Surname         string
	Slug            string
	Email           string
	Password        string
	AffiliateCode   string
	AffiliatePoints int
	IsAdmin         bool
	IsAuthor        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
