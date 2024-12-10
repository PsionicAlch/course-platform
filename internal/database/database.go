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
	GetUsers(term string, level AuthorizationLevel, likedTutorialID, bookmarkedTutorialID string) ([]*models.UserModel, error)
	GetUsersPaginated(term string, level AuthorizationLevel, likedTutorialID, bookmarkedTutorialID string, page, elements uint) ([]*models.UserModel, error)
	GetAllUsers() ([]*models.UserModel, error)
	AddNewUser(name, surname, email, password, token, tokenType, ipAddr string, validUntil time.Time) (*models.UserModel, error)
	NewUser(name, surname, email, password string) error
	NewAdminUser(name, surname, email, password string) error
	GetUserByEmail(email string, level AuthorizationLevel) (*models.UserModel, error)
	GetUserByID(id string, level AuthorizationLevel) (*models.UserModel, error)
	GetUserByToken(token, tokenType string) (*models.UserModel, error)
	GetUserByAffiliateCode(affiliateCode string) (*models.UserModel, error)
	UpdateUserPassword(userId, password string) error
	CountUsers() (uint, error)
	AddAuthorStatus(userId string) error
	RemoveAuthorStatus(userId string) error
	AddAdminStatus(userId string) error
	RemoveAdminStatus(userId string) error

	// Tokens functions.
	AddToken(token, tokenType, userId string, validUntil time.Time) error
	GetToken(token, tokenType string) (*models.TokenModel, error)
	DeleteToken(token, tokenType string) error
	DeleteAllTokens(email, tokenType string) error

	// IP Addresses functions.
	AddIPAddress(userId, ipAddr string) error
	GetUserIpAddresses(userId string) ([]string, error)

	// Tutorials functions.
	AdminGetTutorials(term string, published *bool, authorId *string, likedByUser string, bookmarkedByUser string, keyword string, page, elements uint) ([]*models.TutorialModel, error)
	GetAllTutorials() ([]*models.TutorialModel, error)
	GetTutorials(term string, page, elements int) ([]*models.TutorialModel, error)
	GetTutorialByID(id string) (*models.TutorialModel, error)
	GetTutorialBySlug(slug string) (*models.TutorialModel, error)
	CountTutorials() (uint, error)
	PublishTutorial(tutorialId string) error
	UnpublishTutorial(tutorialId string) error
	UpdateTutorialAuthor(tutorialId, authorId string) error
	CountTutorialsWrittenBy(authorId string) (uint, error)

	// Keywords functions.
	GetKeywords() ([]string, error)
	DeleteAllKeywords() error

	// Tutorials-Keywords functions.
	GetAllKeywordsForTutorial(tutorialId string) ([]string, error)

	// Tutorials-Likes functions.
	GetTutorialsLikedByUser(term, userId string, page, elements uint) ([]*models.TutorialModel, error)
	UserLikedTutorial(userId, slug string) (bool, error)
	UserLikeTutorial(userId, slug string) error
	UserDislikeTutorial(userId, slug string) error
	CountTutorialsLikedByUser(userId string) (uint, error)
	CountTutorialLikes(tutorialId string) (uint, error)

	// Tutorials-Bookmarks functions.
	GetTutorialsBookmarkedByUser(term, userId string, page, elements uint) ([]*models.TutorialModel, error)
	UserBookmarkedTutorial(userId, slug string) (bool, error)
	UserBookmarkTutorial(userId, slug string) error
	UserUnbookmarkTutorial(userId, slug string) error
	CountTutorialsBookmarkedByUser(userId string) (uint, error)
	CountTutorialBookmarks(tutorialId string) (uint, error)

	// Comments functions.
	AdminGetComments(term, tutorialId, userId string, page, elements uint) ([]*models.CommentModel, error)
	GetAllCommentsPaginated(tutorialId string, page, elements int) ([]*models.CommentModel, error)
	GetAllCommentsBySlugPaginated(slug string, page, elements int) ([]*models.CommentModel, error)
	CountCommentsForTutorial(tutorialId string) (uint, error)
	AddCommentBySlug(content, userId, slug string) (*models.CommentModel, error)
	CountComments() (uint, error)
	DeleteComment(commentId string) error

	// Courses functions.
	AdminGetCourses(term string, published *bool, authorId *string, boughtBy, keyword string, page, elements uint) ([]*models.CourseModel, error)
	GetAllCourses() ([]*models.CourseModel, error)
	GetCourses(term string, page, elements int) ([]*models.CourseModel, error)
	GetCourseByFileKey(fileKey string) (*models.CourseModel, error)
	GetCourseBySlug(slug string) (*models.CourseModel, error)
	GetCourseByID(courseId string) (*models.CourseModel, error)
	CountCourses() (uint, error)
	PublishCourse(courseId string) error
	UnpublishCourse(courseId string) error
	UpdateCourseAuthor(tutorialId, authorId string) error

	// Courses Keywords functions.
	GetAllKeywordsForCourse(courseId string) ([]string, error)

	// Chapters functions.
	GetAllChapters() ([]*models.ChapterModel, error)
	GetChapterByFileKey(fileKey string) (*models.ChapterModel, error)
	CountChapters(courseId string) (int, error)

	// Discounts functions.
	GetDiscountsPaginated(term string, active *bool, page, elements uint) ([]*models.DiscountModel, error)
	CountDiscounts() (uint, error)
	AddDiscount(title, description string, discount, uses uint64) error
	GetDiscountByID(discountId string) (*models.DiscountModel, error)
	GetDiscountByCode(discountCode string) (*models.DiscountModel, error)
	ActivateDiscount(discountId string) error
	DeactivateDiscount(discountId string) error

	// Course Purchases functions.
	HasUserPurchasedCourse(userId, courseId string) (bool, error)
	RegisterCoursePurchase(userId, courseId, paymentKey, stripeCheckoutSessionId string, affiliateCode, discountCode sql.NullString, affiliatePointsUsed int64, amountPaid float64, token, tokenType string, validUntil time.Time) error
	CountCoursesWhereDiscountWasUsed(discountCode string) (uint, error)
	CountUsersWhoBoughtCourse(courseId string) (uint, error)
	GetCoursePurchaseByPaymentKey(paymentKey string) (*models.CoursePurchaseModel, error)
	UpdateCoursePurchasePaymentStatus(coursePurchaseId string, status PaymentStatus) error
	GetCoursesBoughtByUser(term, userId string, page, elements uint) ([]*models.CourseModel, error)

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
