package gatekeeper

import "fmt"

type FailedToFindUserByEmail struct {
	msg string
}

func CreateFailedToFindUserByEmail(email, err string) FailedToFindUserByEmail {
	return FailedToFindUserByEmail{
		msg: fmt.Sprintf("failed to find user by %s: %s", email, err),
	}
}

func (err FailedToFindUserByEmail) Error() string {
	return err.msg
}

type UserAlreadyExists struct {
	msg string
}

func CreateUserAlreadyExists(email string) UserAlreadyExists {
	return UserAlreadyExists{
		msg: fmt.Sprintf("a user with the email %s already exists", email),
	}
}

func (err UserAlreadyExists) Error() string {
	return err.msg
}
