package errors

import "fmt"

// FailedToCreateAuthenticationToken represents an inability to create a new authentication token.
type FailedToCreateAuthenticationToken struct {
	msg string
}

// CreateFailedToCreateAuthenticationToken creates a new instance of the FailedToCreateAuthenticationToken error.
func CreateFailedToCreateAuthenticationToken(err string) FailedToCreateAuthenticationToken {
	return FailedToCreateAuthenticationToken{
		msg: fmt.Sprintf("failed to create authentication token: %s", err),
	}
}

// Error returns an error message.
func (err FailedToCreateAuthenticationToken) Error() string {
	return err.msg
}
