package models

import "time"

// ChapterModel is a struct representation of the course_chapters table.
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
