package models

import (
	"database/sql"
	"time"
)

// TutorialModel is a struct representation of the tutorials table.
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
}
