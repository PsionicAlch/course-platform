package errors

import "fmt"

// FailedToMigrate represents a failure in the migration process.
type FailedToMigrate struct {
	msg string
}

// CreateFailedToMigrate creates a new instance of the FailedToMigrate error.
func CreateFailedToMigrate(err string) FailedToMigrate {
	return FailedToMigrate{
		msg: err,
	}
}

// Error returns an error message.
func (err FailedToMigrate) Error() string {
	return fmt.Sprintf("failed to run migrations: %s", err.msg)
}
