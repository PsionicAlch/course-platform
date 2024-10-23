package errors

import "fmt"

// FailedToGenerateToken represents an inability in the system to generate a new token.
type FailedToGenerateToken struct {
	msg string
}

// CreateFailedToGenerateToken creates a new instance of the FailedToGenerateToken error.
func CreateFailedToGenerateToken(err string) FailedToGenerateToken {
	return FailedToGenerateToken{
		msg: fmt.Sprintf("failed to generate new token: %s", err),
	}
}

// Error returns an error message.
func (err FailedToGenerateToken) Error() string {
	return err.msg
}
