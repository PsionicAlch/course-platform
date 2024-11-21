package sqlite_database

func (db *SQLiteDatabase) GetAllKeywordsForTutorial(tutorialId string) ([]string, error) {
	query := `SELECT k.keyword FROM tutorials_keywords tk JOIN keywords k ON tk.keyword_id = k.id WHERE tk.tutorial_id = ?;`

	var keywords []string

	rows, err := db.connection.Query(query, tutorialId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var keyword string

		if err := rows.Scan(&keyword); err != nil {
			return nil, err
		}

		keywords = append(keywords, keyword)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return keywords, nil
}
