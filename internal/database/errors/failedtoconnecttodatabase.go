package errors

import "fmt"

// FailedToConnectToDatabase shows that the program was unable to connect to the database.
type FailedToConnectToDatabase struct {
	msg string
}

// CreateFailedTiConnectToDatabase creates a new instance of the FailedToConnectToDatabase error.
func CreateFailedToConnectToDatabase(err string) FailedToConnectToDatabase {
	return FailedToConnectToDatabase{
		msg: err,
	}
}

// Error returns an error message.
func (err FailedToConnectToDatabase) Error() string {
	return fmt.Sprintf("failed to connect to database: %s", err.msg)
}
