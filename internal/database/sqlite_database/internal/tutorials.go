package internal

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func AddTutorial(dbFacade SqlDbFacade, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey string) error {
	query := `INSERT INTO tutorials (id, title, slug, description, thumbnail_url, banner_url, content, file_checksum, file_key) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	result, err := dbFacade.Exec(query, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey)
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
	query := `UPDATE tutorials SET title = ?, slug = ?, description = ?, thumbnail_url = ?, banner_url = ?, content = ?, published = 0, author_id = ?, file_checksum = ?, file_key = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;`

	result, err := dbFacade.Exec(query, title, slug, description, thumbnailUrl, bannerUrl, content, authorId, fileChecksum, fileKey, id)
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
