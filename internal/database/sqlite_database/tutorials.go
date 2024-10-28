package sqlite_database

import (
	"fmt"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

func (db *SQLiteDatabase) AddKeyword(keyword string) (string, error) {
	return "", nil
}

func (db *SQLiteDatabase) GetAllTutorials() ([]*models.TutorialModel, error) {
	// SQL query to get all tutorials from tutorials table.
	query := `SELECT id, title, slug, description, thumbnail_url, banner_url, content, published, author_id, file_checksum, created_at, updated_at FROM tutorials;`

	// Query the database.
	rows, err := db.connection.Query(query)
	if err != nil {
		db.ErrorLog.Printf("failed to read all the rows of the tutorials table: %s\n", err)
		return nil, err
	}
	defer rows.Close()

	// Create intermediate variable to store all the tutorials.
	var tutorials []*models.TutorialModel

	// Iterate over each row of the database.
	for rows.Next() {
		// Some more intermediate variables to hold each row.
		var tutorial models.TutorialModel
		// SQLite doesn't have booleans so we have to use integers as an intermediate type.
		var publishedInt int

		// Read the data from the table row into the intermediary variables.
		err := rows.Scan(
			&tutorial.ID,
			&tutorial.Title,
			&tutorial.Slug,
			&tutorial.Description,
			&tutorial.ThumbnailURL,
			&tutorial.BannerURL,
			&tutorial.Content,
			&publishedInt,
			&tutorial.AuthorID,
			&tutorial.FileChecksum,
			&tutorial.CreatedAt,
			&tutorial.UpdatedAt,
		)
		if err != nil {
			db.ErrorLog.Printf("Failed to query row of tutorials table: %s\n", err)
			return nil, err
		}

		// Convert the published column to a boolean.
		tutorial.Published = publishedInt == 1

		// Add the individual tutorial to the slice of tutorials.
		tutorials = append(tutorials, &tutorial)
	}

	// Check for any errors after iterating through all the rows.
	if rows.Err() != nil {
		db.ErrorLog.Printf("Found an error after querying all the tutorials table's rows: %s\n", err)
		return nil, err
	}

	return tutorials, nil
}

func (db *SQLiteDatabase) AddNewTutorial(title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum string) error {
	// SQL query to insert new tutorial.
	query := `INSERT INTO tutorials (id, title, slug, description, thumbnail_url, banner_url, content, file_checksum) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

	// Generate new ULID based ID for the tutorial.
	id, err := database.GenerateID()
	if err != nil {
		db.ErrorLog.Printf("Failed to generate new ID for tutorial: %s\n", err)
		return err
	}

	// Execute the SQL query and pass in the required variables.
	_, err = db.connection.Exec(query, id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum)
	if err != nil {
		db.ErrorLog.Printf("Failed to save tutorial \"%s\" to the database: %s\n", title, err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) AddNewTutorialBulk(tutorials []*models.TutorialModel) error {
	// Begin a transaction.
	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("Failed to begin transaction for bulk insert into tutorials table: %s\n", err)
		return err
	}

	// Prepare the SQL query for inserting a new tutorial.
	query := `INSERT INTO tutorials (id, title, slug, description, thumbnail_url, banner_url, content, file_checksum) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		db.ErrorLog.Printf("Failed to prepare statement for bulk insert into tutorials table: %s\n", err)
		return err
	}
	defer stmt.Close() // Ensure the statement is closed after the transaction completes

	// Loop through each tutorial and insert it.
	for _, tutorial := range tutorials {
		// Generate a new ULID-based ID for the tutorial.
		id, err := database.GenerateID()
		if err != nil {
			tx.Rollback()
			db.ErrorLog.Printf("Failed to generate new ID for tutorial: %s\n", err)
			return err
		}

		// Execute the statement for the current tutorial.
		_, err = stmt.Exec(id, tutorial.Title, tutorial.Slug, tutorial.Description, tutorial.ThumbnailURL, tutorial.BannerURL, tutorial.Content, tutorial.FileChecksum)
		if err != nil {
			tx.Rollback()
			db.ErrorLog.Printf("Failed to save tutorial \"%s\" to the database: %s\n", tutorial.Title, err)
			return err
		}
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.ErrorLog.Printf("Failed to commit transaction for bulk insert into tutorials table: %s\n", err)
		return err
	}

	return nil
}

func (db *SQLiteDatabase) UpdateTutorial(id string, tutorial *models.TutorialModel) error {
	query := `UPDATE tutorials SET title = ?, slug = ?, description = ?, thumbnail_url = ?, banner_url = ?, content = ?, published = 0, author_id = null, file_checksum = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;`

	result, err := db.connection.Exec(query, tutorial.Title, tutorial.Slug, tutorial.Description, tutorial.ThumbnailURL, tutorial.BannerURL, tutorial.Content, tutorial.FileChecksum, id)
	if err != nil {
		db.ErrorLog.Printf("Failed update tutorial %s: %s\n", id, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.ErrorLog.Printf("Failed to get the affected rows after updating tutorial %s: %s\n", id, err)
		return err
	}

	if rowsAffected == 0 {
		db.ErrorLog.Printf("Failed to update tutorial %s. 0 rows were affected by this update\n", id)
		return fmt.Errorf("failed to update tutorial %s. 0 rows were affected by this update", id)
	}

	return nil
}

func (db *SQLiteDatabase) UpdateTutorialBulk(tutorials []*models.TutorialModel) error {
	// Start a transaction.
	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("failed to begin transaction for bulk update to tutorials table: %s\n", err)
		return err
	}

	// Prepare the SQL update statement.
	query := `UPDATE tutorials SET title = ?, slug = ?, description = ?, thumbnail_url = ?, banner_url = ?, content = ?, published = 0, author_id = ?, file_checksum = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;`
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		db.ErrorLog.Printf("Failed to prepare statement for bulk update: %s\n", err)
		return err
	}
	defer stmt.Close()

	// Execute the update for each tutorial in the slice.
	for _, tutorial := range tutorials {
		result, err := stmt.Exec(tutorial.Title, tutorial.Slug, tutorial.Description, tutorial.ThumbnailURL, tutorial.BannerURL, tutorial.Content, tutorial.AuthorID, tutorial.FileChecksum, tutorial.ID)
		if err != nil {
			tx.Rollback()
			db.ErrorLog.Printf("Failed to update tutorial %s: %s\n", tutorial.ID, err)
			return err
		}

		// Check rows affected.
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			db.ErrorLog.Printf("Failed to get affected rows after updating tutorial %s: %s\n", tutorial.ID, err)
			return err
		}

		if rowsAffected == 0 {
			tx.Rollback()
			db.ErrorLog.Printf("Failed to update tutorial %s: no rows were affected\n", tutorial.ID)
			return fmt.Errorf("failed to update tutorial %s: no rows were affected", tutorial.ID)
		}
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		db.ErrorLog.Printf("Failed to commit bulk update transaction: %s\n", err)
		return err
	}

	return nil
}
