package database

import "errors"

var (
	ErrNoRowsAffected    = errors.New("no database rows were affected")
	ErrUserAlreadyExists = errors.New("user already exists")
)
