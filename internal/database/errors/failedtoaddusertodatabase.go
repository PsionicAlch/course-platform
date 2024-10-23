package errors

import "fmt"

// FailedToAddUserToDatabase represents the system's inability to add a new user instance to the database.
type FailedToAddUserToDatabase struct {
	msg string
}

// CreateFailedToAddUserToDatabase creates a new instance of the FailedToAddUserToDatabase error.
func CreateFailedToAddUserToDatabase(err string) FailedToAddUserToDatabase {
	return FailedToAddUserToDatabase{
		msg: fmt.Sprintf("failed to add user to the database: %s", err),
	}
}

// Error returns an error message.
func (err FailedToAddUserToDatabase) Error() string {
	return err.msg
}
