package models

import "time"

type UserModel struct {
	ID            string
	Name          string
	Surname       string
	Email         string
	Password      string
	AffiliateCode string
	IsAdmin       bool
	IsAuthor      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type User struct {
	ID            string
	Name          string
	Surname       string
	Email         string
	AffiliateCode string
	IsAdmin       bool
	IsAuthor      bool
}
