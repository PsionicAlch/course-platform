package internal

import "github.com/PsionicAlch/psionicalch-home/internal/database"

func AddCourse(dbFacade SqlDbFacade, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey string) error {
	query := `INSERT INTO courses (id, title, slug, description, thumbnail_url, banner_url, content, file_checksum, file_key) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	results, err := dbFacade.Exec(query, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNoRowsAffected
	}

	return nil
}

func UpdateCourse(dbFacade SqlDbFacade, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey string) error {
	query := `UPDATE courses SET title = ?, slug = ?, description = ?, thumbnail_url = ?, banner_url = ?, content = ?, published = 0, file_checksum = ?, file_key = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;`
	results, err := dbFacade.Exec(query, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey, id)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNoRowsAffected
	}

	return nil
}
