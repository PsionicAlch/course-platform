package models

import (
	"time"
)

// CommentModel is a struct representation of the comments table.
type CommentModel struct {
	ID         string
	Content    string
	UserID     string
	TutorialID string
	CreatedAt  time.Time
}
