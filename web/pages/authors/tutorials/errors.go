package tutorials

import "errors"

var (
	ErrAuthorNotFound    = errors.New("author was not found")
	ErrTutorialsNotFound = errors.New("tutorials were not found")
)
