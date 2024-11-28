package sqlite_database

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/xeonx/timeago"
)

func (db *SQLiteDatabase) CommentSetUser(comment *models.CommentModel) error {
	user, err := db.GetUserByID(comment.UserID, database.All)
	if err != nil {
		db.ErrorLog.Printf("Failed to get user from comment: %s\n", err)
		return err
	}

	comment.User = user

	return nil
}

func (db *SQLiteDatabase) CommentsSetUser(comments []*models.CommentModel) error {
	for _, comment := range comments {
		user, err := db.GetUserByID(comment.UserID, database.All)
		if err != nil {
			db.ErrorLog.Printf("Failed to get user from comment: %s\n", err)
			return err
		}

		comment.User = user
	}

	return nil
}

func (db *SQLiteDatabase) CommentSetTutorial(comment *models.CommentModel) error {
	tutorial, err := db.GetTutorialByID(comment.TutorialID)
	if err != nil {
		db.ErrorLog.Printf("Failed to get comment's (\"%s\") tutorial from the database: %s\n", comment.ID, err)
		return err
	}

	comment.Tutorial = tutorial

	return nil
}

func (db *SQLiteDatabase) CommentsSetTutorial(comments []*models.CommentModel) error {
	for _, comment := range comments {
		tutorial, err := db.GetTutorialByID(comment.TutorialID)
		if err != nil {
			db.ErrorLog.Printf("Failed to get comment's (\"%s\") tutorial from the database: %s\n", comment.ID, err)
			return err
		}

		comment.Tutorial = tutorial
	}

	return nil
}

func (db *SQLiteDatabase) CommentSetTimeAgo(comment *models.CommentModel) {
	comment.TimeAgo = timeago.English.Format(comment.CreatedAt)
}

func (db *SQLiteDatabase) CommentsSetTimeAgo(comments []*models.CommentModel) {
	for _, comment := range comments {
		comment.TimeAgo = timeago.English.Format(comment.CreatedAt)
	}
}
