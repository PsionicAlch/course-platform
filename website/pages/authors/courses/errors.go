package courses

import "errors"

var (
	ErrAuthorNotFound  = errors.New("author not found")
	ErrCoursesNotFound = errors.New("courses not found")
)
