package models

import "time"

type ChapterModel struct {
	ID           string
	Title        string
	Slug         string
	Chapter      int
	Content      string
	CourseID     string
	FileChecksum string
	FileKey      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
