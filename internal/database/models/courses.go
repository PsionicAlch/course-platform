package models

import (
	"database/sql"
	"time"
)

type CourseModel struct {
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

	Keywords []string
}
