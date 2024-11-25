package database

import (
	"database/sql"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

// TODO: Rewrite query functions for tutorials and courses to only show published tutorials and courses.
type Database interface {
	// General database functions.
	Close() error

	// Migration functions.
	MigrateUp() error
	MigrateDown() error
	Rollback(steps int) error

	// Users functions.
	GetAllAdminsPaginated(page, elements int) ([]*models.UserModel, error)
	GetAllUsersPaginated(page, elements int) ([]*models.UserModel, error)
	GetAllAuthorsPaginated(page, elements int) ([]*models.UserModel, error)
	AddNewUser(name, surname, email, password, token, tokenType, ipAddr string, validUntil time.Time) (*models.UserModel, error)
	NewUser(name, surname, email, password string) error
	NewAdminUser(name, surname, email, password string) error
	GetUserByEmail(email string) (*models.UserModel, error)
	GetUserByID(id string) (*models.UserModel, error)
	GetUserByToken(token, tokenType string) (*models.UserModel, error)
	UpdateUserPassword(userId, password string) error

	// Tokens functions.
	AddToken(token, tokenType, userId string, validUntil time.Time) error
	GetToken(token, tokenType string) (*models.TokenModel, error)
	DeleteToken(token, tokenType string) error
	DeleteAllTokens(email, tokenType string) error

	// IP Addresses functions.
	AddIPAddress(userId, ipAddr string) error
	GetUserIpAddresses(userId string) ([]string, error)

	// Tutorials functions.
	GetAllTutorials() ([]*models.TutorialModel, error)
	GetAllTutorialsPaginated(page, elements int) ([]*models.TutorialModel, error)
	SearchTutorialsPaginated(term string, page, elements int) ([]*models.TutorialModel, error)
	GetTutorialBySlug(slug string) (*models.TutorialModel, error)

	// Tutorials-Keywords functions.
	GetAllKeywordsForTutorial(tutorialId string) ([]string, error)

	// Tutorials-Likes functions.
	UserLikedTutorial(userId, slug string) (bool, error)
	UserLikeTutorial(userId, slug string) error
	UserDislikeTutorial(userId, slug string) error

	// Tutorials-Bookmarks functions.
	UserBookmarkedTutorial(userId, slug string) (bool, error)
	UserBookmarkTutorial(userId, slug string) error
	UserUnbookmarkTutorial(userId, slug string) error

	// Comments functions.
	GetAllCommentsPaginated(tutorialId string, page, elements int) ([]*models.CommentModel, error)
	GetAllCommentsBySlugPaginated(slug string, page, elements int) ([]*models.CommentModel, error)
	AddCommentBySlug(content, userId, slug string) (*models.CommentModel, error)

	// Courses functions.
	GetAllCourses() ([]*models.CourseModel, error)
	GetAllCoursesPaginated(page, elements int) ([]*models.CourseModel, error)
	SearchCoursesPaginated(term string, page, elements int) ([]*models.CourseModel, error)
	GetCourseByFileKey(fileKey string) (*models.CourseModel, error)
	GetCourseBySlug(slug string) (*models.CourseModel, error)

	// Chapters functions.
	GetAllChapters() ([]*models.ChapterModel, error)
	GetChapterByFileKey(fileKey string) (*models.ChapterModel, error)
	CountChapters(courseId string) (int, error)

	// Models functions.
	CommentSetUser(comment *models.CommentModel) error
	CommentsSetUser(comments []*models.CommentModel) error
	CommentSetTutorial(comment *models.CommentModel) error
	CommentsSetTutorial(comments []*models.CommentModel) error
	CommentSetTimeAgo(comment *models.CommentModel)
	CommentsSetTimeAgo(comment []*models.CommentModel)

	// Bulk functions.
	PrepareBulkTutorials()
	InsertTutorial(title, slug, description, thumbnailUrl, bannerUrl, content, checksum, fileKey string, keywords []string)
	UpdateTutorial(id, title, slug, description, thumbnailUrl, bannerUrl, content, checksum, fileKey string, keywords []string, authorId sql.NullString)
	RunBulkTutorials() error

	PrepareBulkCourses()
	InsertCourse(title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey string, keywords []string)
	UpdateCourse(id, title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum, fileKey string, keywords []string, authorId sql.NullString)
	InsertChapter(title string, chapter int, content, fileChecksum, fileKey, courseKey string)
	UpdateChapter(id, title string, chapter int, content, fileChecksum, fileKey, courseKey string)
	RunBulkCourses() error
}
