package models

import (
	"database/sql"
	"time"
)

type TutorialModel struct {
	ID           string
	Title        string
	Slug         string
	Description  string
	ThumbnailURL string
	BannerURL    string
	Content      string
	Published    bool
	AuthorID     sql.NullString
	FileChecksum string
	FileKey      string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Keywords []*KeywordModel
}

type Tutorial struct {
	Title        string
	Slug         string
	Description  string
	ThumbnailURL string
	BannerURL    string
	Content      string
	Published    bool
	Author       *User
	Keywords     []Keyword
}
