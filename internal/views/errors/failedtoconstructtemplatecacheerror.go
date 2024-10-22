package errors

import "fmt"

type FailedToConstructTemplateCacheError struct {
	msg string
}

func CreateFailedToConstructTemplateCacheError(err string) FailedToConstructTemplateCacheError {
	return FailedToConstructTemplateCacheError{
		msg: fmt.Sprintf("failed to construct template cache: %s", err),
	}
}

func (err FailedToConstructTemplateCacheError) Error() string {
	return err.msg
}
