package database

import "errors"

var (
	ErrNoRowsAffected        = errors.New("no database rows were affected")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrTutorialAlreadyExists = errors.New("tutorial already exists")
	ErrKeywordAlreadyExists  = errors.New("keyword already exists")
)
