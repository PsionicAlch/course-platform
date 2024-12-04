package sqlite_database

func (db *SQLiteDatabase) GetAllKeywordsForCourse(courseId string) ([]string, error) {
	query := `SELECT k.keyword FROM courses_keywords AS ck JOIN keywords AS k ON ck.keyword_id = k.id WHERE ck.course_id = ?;`

	var keywords []string

	rows, err := db.connection.Query(query, courseId)
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
