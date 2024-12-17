package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

func (db *SQLiteDatabase) AdminGetTutorials(term string, published *bool, authorId *string, likedByUser, bookmarkedByUser, keyword string, page, elements uint) ([]*models.TutorialModel, error) {
	query := `SELECT DISTINCT t.id, t.title, t.slug, t.description, t.thumbnail_url, t.banner_url, t.content, t.published, t.author_id, t.file_checksum, t.file_key, t.created_at, t.updated_at FROM tutorials AS t LEFT JOIN tutorials_likes AS tl ON t.id = tl.tutorial_id LEFT JOIN tutorials_bookmarks AS tb ON t.id = tb.tutorial_id LEFT JOIN tutorials_keywords AS tk ON t.id = tk.tutorial_id LEFT JOIN keywords AS k ON tk.keyword_id = k.id WHERE (LOWER(t.id) LIKE '%' || ? || '%' OR LOWER(t.title) LIKE '%' || ? || '%' OR LOWER(t.slug) LIKE '%' || ? || '%' OR LOWER(t.description) LIKE '%' || ? || '%' OR LOWER(k.keyword) LIKE '%' || ? || '%')`

	args := []any{term, term, term, term, term}

	if published != nil {
		query += " AND t.published = ?"

		var pubInt int
		if *published {
			pubInt = 1
		} else {
			pubInt = 0
		}

		args = append(args, pubInt)
	}

	if authorId != nil {
		if *authorId != "" {
			query += " AND t.author_id = ?"
			args = append(args, *authorId)
		}
	} else {
		query += " AND t.author_id IS NULL"
	}

	if likedByUser != "" {
		query += " AND tl.user_id = ?"
		args = append(args, likedByUser)
	}

	if bookmarkedByUser != "" {
		query += " AND tb.user_id = ?"
		args = append(args, bookmarkedByUser)
	}

	if keyword != "" {
		query += " AND k.keyword LIKE '%' || ? || '%'"
		args = append(args, keyword)
	}

	offset := (page - 1) * elements
	query += " ORDER BY t.created_at DESC, t.updated_at DESC, t.title ASC LIMIT ? OFFSET ?;"
	args = append(args, elements, offset)

	var tutorials []*models.TutorialModel

	rows, err := db.connection.Query(query, args...)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all tutorials from the database: %s\n", err)
		db.ErrorLog.Printf("\nSQL Query Used:\n%s\n", query)

		return nil, err
	}

	for rows.Next() {
		var tutorial models.TutorialModel
		var published int

		if err := rows.Scan(&tutorial.ID, &tutorial.Title, &tutorial.Slug, &tutorial.Description, &tutorial.ThumbnailURL, &tutorial.BannerURL, &tutorial.Content, &published, &tutorial.AuthorID, &tutorial.FileChecksum, &tutorial.FileKey, &tutorial.CreatedAt, &tutorial.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from tutorials table: %s\n", err)
			return nil, err
		}

		tutorial.Published = published == 1

		tutorials = append(tutorials, &tutorial)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to get all tutorials from the database: %s\n", err)
		db.ErrorLog.Printf("\nSQL Query Used:\n%s\n", query)

		return nil, err
	}

	return tutorials, nil
}

func (db *SQLiteDatabase) GetAllTutorials(published *bool) ([]*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM tutorials`
	args := []any{}

	if published != nil {
		query += " WHERE published = ?"

		if *published {
			query += " AND author_id IS NOT NULL"
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}

	query += " ORDER BY updated_at DESC, title ASC;"

	var tutorials []*models.TutorialModel

	rows, err := db.connection.Query(query, args...)
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

func (db *SQLiteDatabase) GetTutorials(term string, authorId string, page, elements int) ([]*models.TutorialModel, error) {
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, file_key, created_at, updated_at FROM tutorials WHERE published = 1 AND author_id IS NOT NULL AND (LOWER(title) LIKE '%' || ? || '%' OR LOWER(description) LIKE '%' || ? || '%')`
	args := []any{term, term}

	if authorId != "" {
		query += " AND author_id = ?"
		args = append(args, authorId)
	}

	offset := (page - 1) * elements
	query += " ORDER BY updated_at DESC, title ASC LIMIT ? OFFSET ?;"
	args = append(args, elements, offset)

	var tutorials []*models.TutorialModel

	rows, err := db.connection.Query(query, args...)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all tutorials (page %d), that match the search term \"%s\", from the database: %s\n", page, term, err)
		return nil, err
	}

	for rows.Next() {
		var tutorial models.TutorialModel
		var published int

		if err := rows.Scan(&tutorial.ID, &tutorial.Title, &tutorial.Slug, &tutorial.Description, &tutorial.ThumbnailURL, &tutorial.BannerURL, &tutorial.Content, &published, &tutorial.AuthorID, &tutorial.FileChecksum, &tutorial.FileKey, &tutorial.CreatedAt, &tutorial.UpdatedAt); err != nil {
			db.ErrorLog.Printf("Failed to read row from tutorials (page %d) table, that match the search term \"%s\": %s\n", page, term, err)
			return nil, err
		}

		if published == 1 {
			tutorial.Published = true
		}

		tutorials = append(tutorials, &tutorial)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to get all tutorials (page %d), that match the search term \"%s\", from the database: %s\n", page, term, err)
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

		db.ErrorLog.Printf("Failed to get tutorial by ID: %s\n", err)
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

		db.ErrorLog.Printf("Failed to get tutorial by slug: %s\n", err)
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
		if err == sql.ErrNoRows {
			return nil, nil
		}

		db.ErrorLog.Printf("Failed to get tutorial by file key from the database: %s\n", err)
		return nil, err
	}

	if publishedInt == 1 {
		tutorial.Published = true
	}

	return &tutorial, nil
}

func (db *SQLiteDatabase) CountTutorials() (uint, error) {
	query := `SELECT COUNT(id) FROM tutorials;`

	var count uint

	row := db.connection.QueryRow(query)
	if err := row.Scan(&count); err != nil {
		db.ErrorLog.Printf("Failed to count the number of tutorials in the database: %s\n", err)
		return 0, err
	}

	return count, nil
}

func (db *SQLiteDatabase) CountTutorialsWrittenBy(authorId string) (uint, error) {
	query := `SELECT COUNT(id) FROM tutorials where author_id = ?;`

	var count uint

	row := db.connection.QueryRow(query, authorId)
	if err := row.Scan(&count); err != nil {
		db.ErrorLog.Printf("Failed to count the number of tutorials written by \"%s\": %s\n", authorId, err)
		return 0, err
	}

	return count, nil
}

func (db *SQLiteDatabase) PublishTutorial(tutorialId string) error {
	query := `UPDATE tutorials SET published = 1 WHERE id = ?;`

	_, err := db.connection.Exec(query, tutorialId)
	if err != nil {
		db.ErrorLog.Printf("Failed to publish tutorial \"%s\": %s\n", tutorialId, err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) UnpublishTutorial(tutorialId string) error {
	query := `UPDATE tutorials SET published = 0 WHERE id = ?;`

	_, err := db.connection.Exec(query, tutorialId)
	if err != nil {
		db.ErrorLog.Printf("Failed to unpublish tutorial \"%s\": %s\n", tutorialId, err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) UpdateTutorialAuthor(tutorialId, authorId string) error {
	query := `UPDATE tutorials SET author_id = ? WHERE id = ?;`

	if authorId == "" {
		authorId = "NULL"
	}

	_, err := db.connection.Exec(query, authorId, tutorialId)
	if err != nil {
		db.ErrorLog.Printf("Failed to update tutorial's (\"%s\") author (\"%s\"): %s\n", tutorialId, authorId, err)
		return err
	}

	return nil
}
