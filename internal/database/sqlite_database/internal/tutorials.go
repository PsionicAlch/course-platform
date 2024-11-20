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
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM tutorials ORDER BY updated_at DESC, title ASC;`

	var tutorials []*models.TutorialModel

	rows, err := dbFacade.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tutorial models.TutorialModel
		var publishedInt int

		published := false

		if err := rows.Scan(&tutorial.ID, &tutorial.Title, &tutorial.Slug, &tutorial.Description, &tutorial.ThumbnailURL, &tutorial.BannerURL, &tutorial.Content, &publishedInt, &tutorial.AuthorID, &tutorial.FileChecksum, &tutorial.FileKey, &tutorial.CreatedAt, &tutorial.UpdatedAt); err != nil {
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

func GetAllTutorialsPaginated(dbFacade SqlDbFacade, page, elements int) ([]*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM tutorials ORDER BY updated_at DESC, title ASC LIMIT ? OFFSET ?;`

	offset := (page - 1) * elements

	var tutorials []*models.TutorialModel

	rows, err := dbFacade.Query(query, elements, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tutorial models.TutorialModel
		var published int

		if err := rows.Scan(&tutorial.ID, &tutorial.Title, &tutorial.Slug, &tutorial.Description, &tutorial.ThumbnailURL, &tutorial.BannerURL, &tutorial.Content, &published, &tutorial.AuthorID, &tutorial.FileChecksum, &tutorial.FileKey, &tutorial.CreatedAt, &tutorial.UpdatedAt); err != nil {
			return nil, err
		}

		if published == 1 {
			tutorial.Published = true
		}

		tutorials = append(tutorials, &tutorial)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tutorials, nil
}

func SearchTutorialsPaginated(dbFacade SqlDbFacade, term string, page, elements int) ([]*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM tutorials WHERE LOWER(title) LIKE '%' || ? || '%' OR LOWER(description) LIKE '%' || ? || '%' ORDER BY updated_at DESC, title ASC LIMIT ? OFFSET ?;`

	offset := (page - 1) * elements

	var tutorials []*models.TutorialModel

	rows, err := dbFacade.Query(query, term, term, elements, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tutorial models.TutorialModel
		var published int

		if err := rows.Scan(&tutorial.ID, &tutorial.Title, &tutorial.Slug, &tutorial.Description, &tutorial.ThumbnailURL, &tutorial.BannerURL, &tutorial.Content, &published, &tutorial.AuthorID, &tutorial.FileChecksum, &tutorial.FileKey, &tutorial.CreatedAt, &tutorial.UpdatedAt); err != nil {
			return nil, err
		}

		if published == 1 {
			tutorial.Published = true
		}

		tutorials = append(tutorials, &tutorial)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tutorials, nil
}

func GetTutorialByID(dbFacade SqlDbFacade, id string) (*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM tutorials WHERE id = ?;`

	var tutorial models.TutorialModel
	var publishedInt int

	row := dbFacade.QueryRow(query, id)
	if err := row.Scan(&tutorial.ID, &tutorial.Title, &tutorial.Slug, &tutorial.Description, &tutorial.ThumbnailURL, &tutorial.BannerURL, &tutorial.Content, &publishedInt, &tutorial.AuthorID, &tutorial.FileChecksum, &tutorial.FileKey, &tutorial.CreatedAt, &tutorial.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if publishedInt == 1 {
		tutorial.Published = true
	}

	return &tutorial, nil
}

func GetTutorialBySlug(dbFacade SqlDbFacade, slug string) (*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM tutorials WHERE slug = ?;`

	var tutorial models.TutorialModel
	var publishedInt int

	row := dbFacade.QueryRow(query, slug)
	if err := row.Scan(&tutorial.ID, &tutorial.Title, &tutorial.Slug, &tutorial.Description, &tutorial.ThumbnailURL, &tutorial.BannerURL, &tutorial.Content, &publishedInt, &tutorial.AuthorID, &tutorial.FileChecksum, &tutorial.FileKey, &tutorial.CreatedAt, &tutorial.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if publishedInt == 1 {
		tutorial.Published = true
	}

	return &tutorial, nil
}

func GetTutorialByFileKey(dbFacade SqlDbFacade, fileKey string) (*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM tutorials WHERE file_key = ?;`

	var tutorial models.TutorialModel
	var publishedInt int

	row := dbFacade.QueryRow(query, fileKey)
	if err := row.Scan(&tutorial.ID, &tutorial.Title, &tutorial.Slug, &tutorial.Description, &tutorial.ThumbnailURL, &tutorial.BannerURL, &tutorial.Content, &publishedInt, &tutorial.AuthorID, &tutorial.FileChecksum, &tutorial.FileKey, &tutorial.CreatedAt, &tutorial.UpdatedAt); err != nil {
		return nil, err
	}

	if publishedInt == 1 {
		tutorial.Published = true
	}

	return &tutorial, nil
}

func AddTutorial(dbFacade SqlDbFacade, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey string) error {
	query := `INSERT INTO tutorials (id, title, slug, description, thumbnail_url, banner_url, content, file_checksum, file_key, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	createdAt := time.Now()

	result, err := dbFacade.Exec(query, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey, createdAt, createdAt)
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

func UpdateTutorial(dbFacade SqlDbFacade, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey string, authorId sql.NullString) error {
	query := `UPDATE tutorials SET title = ?, slug = ?, description = ?, thumbnail_url = ?, banner_url = ?, content = ?, published = 0, author_id = ?, file_checksum = ?, file_key = ?, updated_at = ? WHERE id = ?;`

	updatedAt := time.Now()

	result, err := dbFacade.Exec(query, title, slug, description, thumbnailUrl, bannerUrl, content, authorId, fileChecksum, fileKey, updatedAt, id)
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
