package internal

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
)

func UserBookmarkedTutorial(dbFacade SqlDbFacade, userId, slug string) (bool, error) {
	query := `SELECT t.id FROM tutorials_bookmarks AS tb JOIN tutorials AS t ON tb.tutorial_id = t.id WHERE tb.user_id = ? AND t.slug = ?;`

	var id string

	row := dbFacade.QueryRow(query, userId, slug)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return id != "", nil
}

func UserBookmarkTutorial(dbFacade SqlDbFacade, id, userId, slug string) error {
	query := `INSERT INTO tutorials_bookmarks (id, user_id, tutorial_id) VALUES (?, ?, (SELECT id FROM tutorials WHERE slug = ?));`

	result, err := dbFacade.Exec(query, id, userId, slug)
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

func UserUnbookmarkTutorial(dbFacade SqlDbFacade, userId, slug string) error {
	query := `DELETE FROM tutorials_bookmarks WHERE user_id = ? AND tutorial_id = (SELECT id FROM tutorials WHERE slug = ?);`

	result, err := dbFacade.Exec(query, userId, slug)
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
