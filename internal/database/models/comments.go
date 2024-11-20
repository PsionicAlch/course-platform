package models

import (
	"time"
)

type CommentModel struct {
	ID         string
	Content    string
	UserID     string
	TutorialID string
	CreatedAt  time.Time

	User     *UserModel
	Tutorial *TutorialModel
	TimeAgo  string
}
