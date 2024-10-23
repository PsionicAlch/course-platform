package errors

import "fmt"

// FailedToGenerateID represents a failure to generate a new ID for the database.
type FailedToGenerateID struct {
	msg string
}

// CreateFailedToGenerateID creates a new instance of the FailedToGenerateID error.
func CreateFailedToGenerateID(err string) FailedToGenerateID {
	return FailedToGenerateID{
		msg: fmt.Sprintf("failed to generate new ID: %s", err),
	}
}

// Error returns an error message.
func (err FailedToGenerateID) Error() string {
	return err.msg
}
