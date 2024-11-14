package models

import "time"

type UserModel struct {
	ID            string
	Name          string
	Surname       string
	Email         string
	Password      string
	AffiliateCode string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
