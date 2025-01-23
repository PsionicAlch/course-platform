package sqlite_database

import (
	"database/sql"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

// AdminGetComments gets a paginated list of all comments for the admin panel.
func (db *SQLiteDatabase) AdminGetComments(term, tutorialId, userId string, page, elements uint) ([]*models.CommentModel, error) {
	query := `SELECT c.id, c.content, c.user_id, c.tutorial_id, c.created_at FROM comments AS c WHERE (LOWER(c.id) LIKE '%' || ? || '%' OR LOWER(c.content) LIKE '%' || ? || '%')`

	args := []any{term, term}

	if tutorialId != "" {
		query += " AND c.tutorial_id = ?"
		args = append(args, tutorialId)
	}

	if userId != "" {
		query += " AND c.user_id = ?"
		args = append(args, userId)
	}

	offset := (page - 1) * elements
	query += " ORDER BY c.created_at DESC LIMIT ? OFFSET ?;"
	args = append(args, elements, offset)

	var comments []*models.CommentModel

	rows, err := db.connection.Query(query, args...)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all comments from the database: %s\n", err)
		db.ErrorLog.Printf("\nSQL Query Used:\n%s\n", query)

		return nil, err
	}

	for rows.Next() {
		var comment models.CommentModel

		if err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.TutorialID, &comment.CreatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from comments table: %s\n", err)
			return nil, err
		}

		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to get all comments from the database: %s\n", err)
		db.ErrorLog.Printf("\nSQL Query Used:\n%s\n", query)

		return nil, err
	}

	return comments, nil
}

// GetAllCommentsPaginated gets a paginated list of comments for a given tutorial by ID.
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

// GetAllCommentsBySlugPaginated gets a paginated list of comments for a given tutorial by slug.
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

// CountCommentsForTutorial counts the number of comments a given tutorial has.
func (db *SQLiteDatabase) CountCommentsForTutorial(tutorialId string) (uint, error) {
	query := `SELECT COUNT(id) FROM comments WHERE tutorial_id = ?;`

	var comments uint

	row := db.connection.QueryRow(query, tutorialId)
	if err := row.Scan(&comments); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		db.ErrorLog.Printf("Failed to count all comments related to tutorial \"%s\": %s\n", tutorialId, err)
		return 0, err
	}

	return comments, nil
}

// AddCommentBySlug adds a comment to a tutorial by the tutorial slug.
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

// CountComments counts the number of comments in the database.
func (db *SQLiteDatabase) CountComments() (uint, error) {
	query := `SELECT COUNT(id) FROM comments;`

	var count uint

	row := db.connection.QueryRow(query)
	if err := row.Scan(&count); err != nil {
		db.ErrorLog.Printf("Failed to count all comments in the database: %s\n", err)
		return 0, err
	}

	return count, nil
}

// DeleteComment deletes a comment by it's ID.
func (db *SQLiteDatabase) DeleteComment(commentId string) error {
	query := `DELETE FROM comments WHERE id = ?;`

	result, err := db.connection.Exec(query, commentId)
	if err != nil {
		db.ErrorLog.Printf("Failed to delete comment (\"%s\"): %s\n", commentId, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to get the rows affected after deleting comment: %s\n", err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Printf("No rows were affected after deleting comment (\"%s\")\n", commentId)
		return database.ErrNoRowsAffected
	}

	return nil
}
