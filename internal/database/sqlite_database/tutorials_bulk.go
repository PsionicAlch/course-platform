package sqlite_database

import (
	"database/sql"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database/internal"
)

type intermediate_tutorial struct {
	ID           string
	Title        string
	Slug         string
	Description  string
	ThumbnailURL string
	BannerURL    string
	Content      string
	Checksum     string
	FileKey      string
	Keywords     []string
	AuthorID     sql.NullString
}

var tutorialsToInsert []*intermediate_tutorial
var tutorialsToUpdate []*intermediate_tutorial

func (db *SQLiteDatabase) PrepareBulkTutorials() {
	tutorialsToInsert = []*intermediate_tutorial{}
	tutorialsToUpdate = []*intermediate_tutorial{}
}

func (db *SQLiteDatabase) InsertTutorial(title, slug, description, thumbnailUrl, bannerUrl, content, checksum, fileKey string, keywords []string) {
	tutorialsToInsert = append(tutorialsToInsert, &intermediate_tutorial{
		ID:           "",
		Title:        title,
		Slug:         slug,
		Description:  description,
		ThumbnailURL: thumbnailUrl,
		BannerURL:    bannerUrl,
		Content:      content,
		Checksum:     checksum,
		FileKey:      fileKey,
		Keywords:     keywords,
	})
}

func (db *SQLiteDatabase) UpdateTutorial(id, title, slug, description, thumbnailUrl, bannerUrl, content, checksum, fileKey string, keywords []string, authorId sql.NullString) {
	tutorialsToUpdate = append(tutorialsToUpdate, &intermediate_tutorial{
		ID:           id,
		Title:        title,
		Slug:         slug,
		Description:  description,
		ThumbnailURL: thumbnailUrl,
		BannerURL:    bannerUrl,
		Content:      content,
		Checksum:     checksum,
		FileKey:      fileKey,
		Keywords:     keywords,
		AuthorID:     authorId,
	})
}

func (db *SQLiteDatabase) RunBulkTutorials() error {
	tx, err := db.connection.Begin()
	if err != nil {
		db.ErrorLog.Printf("Failed to create new database transaction: %s\n", err)
		return err
	}

	if err := AddTutorials(tx, tutorialsToInsert); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to bulk insert tutorials: %s\n", err)
		return err
	}

	if err := UpdateTutorials(tx, tutorialsToUpdate); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to bulk update tutorials: %s\n", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			db.ErrorLog.Printf("Failed to rollback changes: %s\n", err)
		}

		db.ErrorLog.Printf("Failed to commit database transaction: %s\n", err)
		return err
	}

	return nil
}

func AddTutorials(tx *sql.Tx, tutorials []*intermediate_tutorial) error {
	for _, tutorial := range tutorials {
		id, err := database.GenerateID()
		if err != nil {
			return err
		}

		if err := internal.AddTutorial(tx, id, tutorial.Title, tutorial.Slug, tutorial.Description, tutorial.ThumbnailURL, tutorial.BannerURL, tutorial.Content, tutorial.Checksum, tutorial.FileKey); err != nil {
			return err
		}

		if err := AddKeywordsToTutorial(tx, id, tutorial.Keywords); err != nil {
			return err
		}
	}

	return nil
}

func UpdateTutorials(tx *sql.Tx, tutorials []*intermediate_tutorial) error {
	for _, tutorial := range tutorials {
		if err := internal.UpdateTutorial(tx, tutorial.ID, tutorial.Title, tutorial.Slug, tutorial.Description, tutorial.ThumbnailURL, tutorial.BannerURL, tutorial.Content, tutorial.Checksum, tutorial.FileKey, tutorial.AuthorID); err != nil {
			return err
		}

		if err := internal.DeleteAllKeywordsFromTutorial(tx, tutorial.ID); err != nil {
			return err
		}

		if err := AddKeywordsToTutorial(tx, tutorial.ID, tutorial.Keywords); err != nil {
			return err
		}
	}

	return nil
}

func AddKeywordsToTutorial(tx *sql.Tx, tutorialId string, keywords []string) error {
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

		tutorialKeywordId, err := database.GenerateID()
		if err != nil {
			return err
		}

		if err := internal.AddKeywordToTutorial(tx, tutorialKeywordId, keywordId, tutorialId); err != nil {
			return err
		}
	}

	return nil
}
