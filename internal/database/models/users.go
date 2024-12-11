package models

import "time"

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
