package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

func (db *SQLiteDatabase) GetAllTutorials() ([]*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM tutorials ORDER BY updated_at DESC, title ASC;`

	var tutorials []*models.TutorialModel

	rows, err := db.connection.Query(query)
	if err != nil {
		db.ErrorLog.Printf("Failed to query database for all tutorials: %s\n", err)
		return nil, err
	}

	for rows.Next() {
		var tutorial models.TutorialModel
		var publishedInt int

		published := false

		if err := rows.Scan(&tutorial.ID, &tutorial.Title, &tutorial.Slug, &tutorial.Description, &tutorial.ThumbnailURL, &tutorial.BannerURL, &tutorial.Content, &publishedInt, &tutorial.AuthorID, &tutorial.FileChecksum, &tutorial.FileKey, &tutorial.CreatedAt, &tutorial.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to scan tutorials row from the database: %s\n", err)
			return nil, err
		}

		if publishedInt == 1 {
			published = true
		}

		tutorial.Published = published

		tutorials = append(tutorials, &tutorial)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Found an error after scanning all tutorials rows: %s\n", err)
		return nil, err
	}

	return tutorials, err
}

func (db *SQLiteDatabase) GetAllTutorialsPaginated(page, elements int) ([]*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM tutorials ORDER BY updated_at DESC, title ASC LIMIT ? OFFSET ?;`

	offset := (page - 1) * elements

	var tutorials []*models.TutorialModel

	rows, err := db.connection.Query(query, elements, offset)
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

func (db *SQLiteDatabase) SearchTutorialsPaginated(term string, page, elements int) ([]*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM tutorials WHERE LOWER(title) LIKE '%' || ? || '%' OR LOWER(description) LIKE '%' || ? || '%' ORDER BY updated_at DESC, title ASC LIMIT ? OFFSET ?;`

	offset := (page - 1) * elements

	var tutorials []*models.TutorialModel

	rows, err := db.connection.Query(query, term, term, elements, offset)
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

func (db *SQLiteDatabase) GetTutorialByID(id string) (*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM tutorials WHERE id = ?;`

	var tutorial models.TutorialModel
	var publishedInt int

	row := db.connection.QueryRow(query, id)
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

func (db *SQLiteDatabase) GetTutorialBySlug(slug string) (*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM tutorials WHERE slug = ?;`

	var tutorial models.TutorialModel
	var publishedInt int

	row := db.connection.QueryRow(query, slug)
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

func (db *SQLiteDatabase) GetTutorialByFileKey(fileKey string) (*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM tutorials WHERE file_key = ?;`

	var tutorial models.TutorialModel
	var publishedInt int

	row := db.connection.QueryRow(query, fileKey)
	if err := row.Scan(&tutorial.ID, &tutorial.Title, &tutorial.Slug, &tutorial.Description, &tutorial.ThumbnailURL, &tutorial.BannerURL, &tutorial.Content, &publishedInt, &tutorial.AuthorID, &tutorial.FileChecksum, &tutorial.FileKey, &tutorial.CreatedAt, &tutorial.UpdatedAt); err != nil {
		db.ErrorLog.Printf("Failed to get tutorial by file key from the database: %s\n", err)
		return nil, err
	}

	if publishedInt == 1 {
		tutorial.Published = true
	}

	return &tutorial, nil
}
