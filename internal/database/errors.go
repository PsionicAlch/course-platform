package database

import "errors"

var (
	ErrNoRowsAffected              = errors.New("no database rows were affected")
	ErrUserAlreadyExists           = errors.New("user already exists")
	ErrTutorialAlreadyExists       = errors.New("tutorial already exists")
	ErrKeywordAlreadyExists        = errors.New("keyword already exists")
	ErrTokenAlreadyExists          = errors.New("token of that type already exists")
	ErrChapterAlreadyExists        = errors.New("that chapter already exists")
	ErrCourseAlreadyOwned          = errors.New("user has already purchased this course")
	ErrInsufficientAffiliatePoints = errors.New("user does not have enough affiliate points")
)
