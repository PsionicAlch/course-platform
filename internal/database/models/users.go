package models

import "time"

type UserModel struct {
	ID              string
	Name            string
	Surname         string
	Email           string
	Password        string
	AffiliateCode   string
	AffiliatePoints uint
	IsAdmin         bool
	IsAuthor        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
