package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/course-platform/internal/database"
	"github.com/PsionicAlch/course-platform/internal/database/sqlite_database/internal"
)

type intermediate_course struct {
	ID           string
	Title        string
	Slug         string
	Description  string
	ThumbnailURL string
	BannerURL    string
	Content      string
	FileChecksum string
	FileKey      string
	Keywords     []string
	AuthorID     sql.NullString
}

type intermediate_chapter struct {
	ID           string
	Title        string
	Slug         string
	Chapter      int
	Content      string
	FileChecksum string
	FileKey      string
	CourseKey    string
}

var coursesToInsert []*intermediate_course
var coursesToUpdate []*intermediate_course
var chaptersToInsert []*intermediate_chapter
var chaptersToUpdate []*intermediate_chapter

func (db *SQLiteDatabase) PrepareBulkCourses() {
	coursesToInsert = []*intermediate_course{}
	coursesToUpdate = []*intermediate_course{}
	chaptersToInsert = []*intermediate_chapter{}
	chaptersToUpdate = []*intermediate_chapter{}
}

func (db *SQLiteDatabase) InsertCourse(title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey string, keywords []string) {
	coursesToInsert = append(coursesToInsert, &intermediate_course{
		Title:        title,
		Slug:         slug,
		Description:  description,
		ThumbnailURL: thumbnailUrl,
		BannerURL:    bannerUrl,
		Content:      content,
		FileChecksum: fileChecksum,
		FileKey:      fileKey,
		Keywords:     keywords,
	})
}

func (db *SQLiteDatabase) UpdateCourse(id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey string, keywords []string, authorId sql.NullString) {
	coursesToUpdate = append(coursesToUpdate, &intermediate_course{
		ID:           id,
		Title:        title,
		Slug:         slug,
		Description:  description,
		ThumbnailURL: thumbnailUrl,
		BannerURL:    bannerUrl,
		Content:      content,
		FileChecksum: fileChecksum,
		FileKey:      fileKey,
		Keywords:     keywords,
		AuthorID:     authorId,
	})
}

func (db *SQLiteDatabase) InsertChapter(title, slug string, chapter int, content, fileChecksum, fileKey, courseKey string) {
	chaptersToInsert = append(chaptersToInsert, &intermediate_chapter{
		Title:        title,
		Slug:         slug,
		Chapter:      chapter,
		Content:      content,
		FileChecksum: fileChecksum,
		FileKey:      fileKey,
		CourseKey:    courseKey,
	})
}

func (db *SQLiteDatabase) UpdateChapter(id, title, slug string, chapter int, content, fileChecksum, fileKey, courseKey string) {
	chaptersToUpdate = append(chaptersToUpdate, &intermediate_chapter{
		ID:           id,
		Title:        title,
		Slug:         slug,
		Chapter:      chapter,
		Content:      content,
		FileChecksum: fileChecksum,
		FileKey:      fileKey,
		CourseKey:    courseKey,
	})
}

func (db *SQLiteDatabase) RunBulkCourses() error {
	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("Failed to start new database transaction for bulk parsing courses: %s\n", err)
		return err
	}

	if err := AddCourses(tx, coursesToInsert); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after an error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to bulk insert courses: %s\n", err)
		return err
	}

	if err := UpdateCourses(tx, coursesToUpdate); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after an error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to bulk update courses: %s\n", err)
		return err
	}

	if err := AddChapters(tx, chaptersToInsert); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after an error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to bulk insert chapters: %s\n", err)
		return err
	}

	if err := UpdateChapters(tx, chaptersToUpdate); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after an error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to bulk update chapters: %s\n", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes after an error occurred: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to commit bulk courses changes to the database: %s\n", err)
		return err
	}

	return nil
}

func AddCourses(tx *sql.Tx, courses []*intermediate_course) error {
	for _, course := range courses {
		id, err := database.GenerateID()
		if err != nil {
			return err
		}

		if err := internal.AddCourse(tx, id, course.Title, course.Slug, course.Description, course.ThumbnailURL, course.BannerURL, course.Content, course.FileChecksum, course.FileKey); err != nil {
			return err
		}

		if err := AddKeywordsToCourse(tx, id, course.Keywords); err != nil {
			return err
		}
	}

	return nil
}

func AddChapters(tx *sql.Tx, chapters []*intermediate_chapter) error {
	for _, chapter := range chapters {
		id, err := database.GenerateID()
		if err != nil {
			return err
		}

		if err := internal.AddChapter(tx, id, chapter.Title, chapter.Slug, chapter.Chapter, chapter.Content, chapter.FileChecksum, chapter.FileKey, chapter.CourseKey); err != nil {
			return err
		}
	}

	return nil
}

func UpdateCourses(tx *sql.Tx, courses []*intermediate_course) error {
	for _, course := range courses {
		if err := internal.UpdateCourse(tx, course.ID, course.Title, course.Slug, course.Description, course.ThumbnailURL, course.BannerURL, course.Content, course.FileChecksum, course.FileKey); err != nil {
			return err
		}

		if err := internal.DeleteAllKeywordsFromCourses(tx, course.ID); err != nil {
			return err
		}

		if err := AddKeywordsToCourse(tx, course.ID, course.Keywords); err != nil {
			return err
		}
	}

	return nil
}

func UpdateChapters(tx *sql.Tx, chapters []*intermediate_chapter) error {
	for _, chapter := range chapters {
		if err := internal.UpdateChapter(tx, chapter.ID, chapter.Title, chapter.Slug, chapter.Chapter, chapter.Content, chapter.FileChecksum, chapter.FileKey, chapter.CourseKey); err != nil {
			return err
		}
	}

	return nil
}

func AddKeywordsToCourse(tx *sql.Tx, courseId string, keywords []string) error {
	for _, keyword := range keywords {
		keywordId, err := database.GenerateID()
		if err != nil {
			return err
		}

		if err := internal.AddKeyword(tx, keywordId, keyword); err != nil {
			if err == database.ErrKeywordAlreadyExists {
				keywordModel, err := internal.GetKeywordByKeyword(tx, keyword)
				if err != nil {
					return err
				}

				keywordId = keywordModel.ID
			} else {
				return err
			}
		}

		courseKeywordId, err := database.GenerateID()
		if err != nil {
			return err
		}

		if err := internal.AddKeywordToCourse(tx, courseKeywordId, keywordId, courseId); err != nil {
			return err
		}
	}

	return nil
}
