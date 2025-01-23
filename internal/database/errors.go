package database

import "errors"

var (
	// ErrNoRowsAffected indicates that no database rows were affected.
	ErrNoRowsAffected = errors.New("no database rows were affected")

	// ErrUserAlreadyExists indicates that the user is already in the database.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrTutorialAlreadyExists indicates that the tutorial is already in the database.
	ErrTutorialAlreadyExists = errors.New("tutorial already exists")

	// ErrKeywordAlreadyExists indicates that the keyword already exists in the database.
	ErrKeywordAlreadyExists = errors.New("keyword already exists")

	// ErrTokenAlreadyExists indicates that the token and token type combination already exists.
	ErrTokenAlreadyExists = errors.New("token of that type already exists")

	// ErrChapterAlreadyExists indicates that the chapter has already been loaded into the database.
	ErrChapterAlreadyExists = errors.New("that chapter already exists")

	// ErrCourseAlreadyOwned indicates that the user already owns the course.
	ErrCourseAlreadyOwned = errors.New("user has already purchased this course")

	// ErrInsufficientAffiliatePoints indicates that there isn't enough affiliate points in the user's profile.
	ErrInsufficientAffiliatePoints = errors.New("user does not have enough affiliate points")
)
