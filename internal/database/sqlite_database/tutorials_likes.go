package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/course-platform/internal/database"
	"github.com/PsionicAlch/course-platform/internal/database/models"
)

func (db *SQLiteDatabase) GetTutorialsLikedByUser(term, userId string, page, elements uint) ([]*models.TutorialModel, error) {
	query := `SELECT t.id, t.title, t.slug, t.description, t.thumbnail_url, t.banner_url, t.content, t.published, t.author_id, t.file_checksum, t.file_key, t.created_at, t.updated_at FROM tutorials_likes AS tl JOIN tutorials AS t ON tl.tutorial_id = t.id WHERE tl.user_id = ? AND t.published = 1`
	args := []any{userId}

	if term != "" {
		query += " AND (LOWER(t.title) LIKE '%' || ? || '%' OR LOWER(t.slug) LIKE '%' || ? || '%' OR LOWER(t.description) LIKE '%' || ? || '%')"
		args = append(args, term, term, term)
	}

	offset := (page - 1) * elements
	query += " ORDER BY tl.created_at DESC LIMIT ? OFFSET ?;"
	args = append(args, elements, offset)

	var tutorials []*models.TutorialModel

	rows, err := db.connection.Query(query, args...)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all tutorials liked by user (\"%s\"): %s\n", userId, err)
		return nil, err
	}

	for rows.Next() {
		var tutorial models.TutorialModel
		var published int

		if err := rows.Scan(&tutorial.ID, &tutorial.Title, &tutorial.Slug, &tutorial.Description, &tutorial.ThumbnailURL, &tutorial.BannerURL, &tutorial.Content, &published, &tutorial.AuthorID, &tutorial.FileChecksum, &tutorial.FileKey, &tutorial.CreatedAt, &tutorial.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read tutorial from the database: %s\n", err)
			return nil, err
		}

		tutorial.Published = published == 1

		tutorials = append(tutorials, &tutorial)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to get all tutorials liked by user (\"%s\"): %s\n", userId, err)
		return nil, err
	}

	return tutorials, nil
}

func (db *SQLiteDatabase) UserLikedTutorial(userId, slug string) (bool, error) {
	query := `SELECT t.id FROM tutorials_likes AS tl JOIN tutorials AS t ON tl.tutorial_id = t.id WHERE tl.user_id = ? AND t.slug = ?;`

	var id string

	row := db.connection.QueryRow(query, userId, slug)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		db.ErrorLog.Printf("Failed to query row from tutorials_likes table: %s\n", err)
		return false, err
	}

	return id != "", nil
}

func (db *SQLiteDatabase) UserLikeTutorial(userId, slug string) error {
	query := `INSERT INTO tutorials_likes (id, user_id, tutorial_id) VALUES (?, ?, (SELECT id FROM tutorials WHERE slug = ?));`

	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate ID for new tutorials_likes database row: %s\n", err)
		return err
	}

	result, err := db.connection.Exec(query, id, userId, slug)
	if err != nil {
		db.ErrorLog.Printf("Failed to add new row to tutorials_likes table: %s\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to query the database for the number of rows affected after adding new row to tutorials_likes table: %s\n", err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Println("No rows were affected after adding a new row to tutorials_likes table")
		return database.ErrNoRowsAffected
	}

	return nil
}

func (db *SQLiteDatabase) UserDislikeTutorial(userId, slug string) error {
	query := `DELETE FROM tutorials_likes WHERE user_id = ? AND tutorial_id = (SELECT id FROM tutorials WHERE slug = ?);`

	result, err := db.connection.Exec(query, userId, slug)
	if err != nil {
		db.ErrorLog.Printf("Failed to delete row from tutorials_likes table: %s\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to query the number of rows affected after deleting row from tutorials_likes table: %s\n", err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Println("No rows were affected after adding new row to tutorials_likes table")
		return database.ErrNoRowsAffected
	}

	return nil
}

func (db *SQLiteDatabase) TutorialsLikedByUser(userId string) ([]*models.TutorialModel, error) {
	query := `SELECT t.id, t.title, t.slug, t.description, t.thumbnail_url, t.banner_url, t.content, t.published, t.author_id, t.file_checksum, t.file_key, t.created_at, t.updated_at FROM tutorials_likes AS tl JOIN tutorials AS t ON tl.tutorial_id = t.id WHERE tl.user_id = ?;`

	var tutorials []*models.TutorialModel

	rows, err := db.connection.Query(query, userId)
	if err != nil {
		db.ErrorLog.Printf("Failed to query database for all tutorials liked by user \"%s\": %s\n", userId, err)
		return nil, err
	}

	for rows.Next() {
		var tutorial models.TutorialModel
		var published int

		if err := rows.Scan(&tutorial.AuthorID, &tutorial.Title, &tutorial.Slug, &tutorial.Description, &tutorial.ThumbnailURL, &tutorial.BannerURL, &tutorial.Content, &published, &tutorial.AuthorID, &tutorial.FileChecksum, &tutorial.FileKey, &tutorial.CreatedAt, &tutorial.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from tutorials_likes table joined with tutorials table: %s\n", err)
			return nil, err
		}

		tutorials = append(tutorials, &tutorial)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to query database for all tutorials liked by user \"%s\": %s\n", userId, err)
		return nil, err
	}

	return tutorials, nil
}

func (db *SQLiteDatabase) CountTutorialsLikedByUser(userId string) (uint, error) {
	query := `SELECT COUNT(id) FROM tutorials_likes WHERE user_id = ?;`

	var count uint

	row := db.connection.QueryRow(query, userId)
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		db.ErrorLog.Printf("Failed to count the number of tutorials liked by user \"%s\": %s\n", userId, err)
		return 0, err
	}

	return count, nil
}

func (db *SQLiteDatabase) CountTutorialLikes(tutorialId string) (uint, error) {
	query := `SELECT COUNT(id) FROM tutorials_likes WHERE tutorial_id = ?;`

	var count uint

	row := db.connection.QueryRow(query, tutorialId)
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		db.ErrorLog.Printf("Failed to count the number of likes the tutorial \"%s\" has: %s\n", tutorialId, err)
		return 0, err
	}

	return count, nil
}
