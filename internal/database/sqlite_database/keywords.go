package sqlite_database

func (db *SQLiteDatabase) GetKeywords() ([]string, error) {
	query := `SELECT keyword FROM keywords;`

	var keywords []string

	rows, err := db.connection.Query(query)
	if err != nil {
		db.ErrorLog.Printf("Failed to get all keywords from the database: %s\n", err)
		return nil, err
	}

	for rows.Next() {
		var keyword string

		if err := rows.Scan(&keyword); err != nil {
			db.ErrorLog.Printf("Failed to read keyword from the database: %s\n", err)
			return nil, err
		}

		keywords = append(keywords, keyword)
	}

	if err := rows.Err(); err != nil {
		db.ErrorLog.Printf("Failed to get all keywords from the database: %s\n", err)
		return nil, err
	}

	return keywords, nil
}
