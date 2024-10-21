package authentication

import "fmt"

// FailedToSignUserIn serves to indicate that the authentication package was unable to sign user in.
type FailedToSignUserIn struct {
	msg string
}

// CreateFailedToSignUserIn creates a new instance of the FailedToSignUserIn error.
func CreateFailedToSignUserIn(err string) FailedToSignUserIn {
	return FailedToSignUserIn{
		msg: fmt.Sprintf("failed to sign user in: %s", err),
	}
}

func (err FailedToSignUserIn) Error() string {
	return err.msg
}
