package sqlite_database

import (
	"database/sql"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

func (db *SQLiteDatabase) GetAllComments(tutorialId string) ([]*models.CommentModel, error) {
	query := `SELECT id, content, user_id, tutorial_id, created_at FROM comments WHERE tutorial_id = ?;`

	var comments []*models.CommentModel

	rows, err := db.connection.Query(query, tutorialId)
	if err != nil {
		if err == sql.ErrNoRows {
			return comments, nil
		}

		db.ErrorLog.Printf("Failed to get all comments from the database: %s\n", err)
		return nil, err
	}

	for rows.Next() {
		var comment models.CommentModel

		if err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.TutorialID, &comment.CreatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read individual row from comments table: %s\n", err)
			return nil, err
		}

		comments = append(comments, &comment)
	}

	if rows.Err() != nil {
		db.ErrorLog.Printf("Got an error after reading all rows from comments table: %s\n", err)
		return nil, err
	}

	return comments, nil
}

func (db *SQLiteDatabase) GetAllCommentsPaginated(tutorialId string, page, elements int) ([]*models.CommentModel, error) {
	query := `SELECT id, content, user_id, tutorial_id, created_at FROM comments WHERE tutorial_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;`

	offset := (page - 1) * elements

	var comments []*models.CommentModel

	rows, err := db.connection.Query(query, tutorialId, elements, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return comments, nil
		}

		db.ErrorLog.Printf("Failed to get comments from the database: %s\n", err)
		return nil, err
	}

	for rows.Next() {
		var comment models.CommentModel

		if err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.TutorialID, &comment.CreatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read individual row from comments table: %s\n", err)
			return nil, err
		}

		comments = append(comments, &comment)
	}

	if rows.Err() != nil {
		db.ErrorLog.Printf("Got an error after reading rows from comments table: %s\n", err)
		return nil, err
	}

	return comments, nil
}

func (db *SQLiteDatabase) GetAllCommentsBySlugPaginated(slug string, page, elements int) ([]*models.CommentModel, error) {
	query := `SELECT id, content, user_id, tutorial_id, created_at FROM comments WHERE tutorial_id = (SELECT id FROM tutorials WHERE slug = ?) ORDER BY created_at DESC LIMIT ? OFFSET ?;`

	offset := (page - 1) * elements

	var comments []*models.CommentModel

	rows, err := db.connection.Query(query, slug, elements, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return comments, nil
		}

		db.ErrorLog.Printf("Failed to get comments from the database: %s\n", err)
		return nil, err
	}

	for rows.Next() {
		var comment models.CommentModel

		if err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.TutorialID, &comment.CreatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read individual row from comments table: %s\n", err)
			return nil, err
		}

		comments = append(comments, &comment)
	}

	if rows.Err() != nil {
		db.ErrorLog.Printf("Got an error after reading rows from comments table: %s\n", err)
		return nil, err
	}

	return comments, nil
}

func (db *SQLiteDatabase) AddCommentBySlug(content, userId, slug string) (*models.CommentModel, error) {
	query := `INSERT INTO comments (id, content, user_id, tutorial_id, created_at) VALUES (?, ?, ?, (SELECT id FROM tutorials WHERE slug = ?), ?);`

	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate ID for new comment: %s\n", err)
		return nil, err
	}

	comment := models.CommentModel{
		ID:         id,
		Content:    content,
		UserID:     userId,
		TutorialID: "",
		CreatedAt:  time.Now(),
	}

	result, err := db.connection.Exec(query, comment.ID, comment.Content, comment.UserID, slug, comment.CreatedAt)
	if err != nil {
		db.ErrorLog.Printf("Failed to insert new comment in comments table: %s\n", err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to read the rows affected after inserting new comment in comments table: %s\n", err)
		return nil, err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Println("No rows were affected after inserting new comment into comments table")
		return nil, database.ErrNoRowsAffected
	}

	return &comment, nil
}
