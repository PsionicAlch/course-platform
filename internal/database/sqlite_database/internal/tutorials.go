package internal

import (
	"database/sql"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func GetAllTutorials(dbFacade SqlDbFacade) ([]*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, created_at, updated_at FROM tutorials ORDER BY updated_at DESC, title ASC;`

	var tutorials []*models.TutorialModel

	rows, err := dbFacade.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tutorial models.TutorialModel
		var publishedInt int

		published := false

		if err := rows.Scan(&tutorial.ID, &tutorial.Title, &tutorial.Slug, &tutorial.Description, &tutorial.ThumbnailURL, &tutorial.BannerURL, &tutorial.Content, &publishedInt, &tutorial.AuthorID, &tutorial.FileChecksum, &tutorial.CreatedAt, &tutorial.UpdatedAt); err != nil {
			return nil, err
		}

		if publishedInt == 1 {
			published = true
		}

		tutorial.Published = published

		tutorials = append(tutorials, &tutorial)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tutorials, err
}

func GetTutorialBySlug(dbFacade SqlDbFacade, slug string) (*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, created_at, updated_at FROM tutorials WHERE slug = ?;`

	var tutorial models.TutorialModel
	var publishedInt int

	row := dbFacade.QueryRow(query, slug)
	if err := row.Scan(&tutorial.ID, &tutorial.Title, &tutorial.Slug, &tutorial.Description, &tutorial.ThumbnailURL, &tutorial.BannerURL, &tutorial.Content, &publishedInt, &tutorial.AuthorID, &tutorial.FileChecksum, &tutorial.CreatedAt, &tutorial.UpdatedAt); err != nil {
		return nil, err
	}

	if publishedInt == 1 {
		tutorial.Published = true
	}

	return &tutorial, nil
}

func AddTutorial(dbFacade SqlDbFacade, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum string) error {
	query := `INSERT INTO tutorials (id, title, slug, description, thumbnail_url, banner_url, content, file_checksum, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	createdAt := time.Now()

	result, err := dbFacade.Exec(query, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, createdAt, createdAt)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return database.ErrTutorialAlreadyExists
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNoRowsAffected
	}

	return nil
}

func UpdateTutorial(dbFacade SqlDbFacade, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum string, authorId sql.NullString) error {
	query := `UPDATE tutorials SET title = ?, slug = ?, description = ?, thumbnail_url = ?, banner_url = ?, content = ?, published = 0, author_id = ?, file_checksum = ?, updated_at = ? WHERE id = ?;`

	updatedAt := time.Now()

	result, err := dbFacade.Exec(query, title, slug, description, thumbnailUrl, bannerUrl, content, authorId, fileChecksum, updatedAt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNoRowsAffected
	}

	return nil
}
