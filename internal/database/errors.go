package database

import "fmt"

// TODO: Cleanup errors!

type FailedToConnectToDatabase struct {
	msg string
}

func CreateFailedConnectToDatabase(err string) FailedToConnectToDatabase {
	return FailedToConnectToDatabase{
		msg: err,
	}
}

func (err FailedToConnectToDatabase) Error() string {
	return fmt.Sprintf("failed to connect to database: %s", err.msg)
}

type FailedToMigrate struct {
	msg string
}

func CreateFailedToMigrate(err string) FailedToMigrate {
	return FailedToMigrate{
		msg: err,
	}
}

func (err FailedToMigrate) Error() string {
	return fmt.Sprintf("failed to run migrations: %s", err.msg)
}

// UserAlreadyExists serves to indicate that the user's email address already exists.
type UserAlreadyExists struct{}

func (err UserAlreadyExists) Error() string {
	return "a user with that email already exists"
}

type FailedToGenerateID struct {
	msg string
}

func CreateFailedToGenerateID(err string) FailedToGenerateID {
	return FailedToGenerateID{
		msg: fmt.Sprintf("failed to generate new ID: %s", err),
	}
}

func (err FailedToGenerateID) Error() string {
	return err.msg
}

type FailedToAddUserToDatabase struct {
	msg string
}

func CreateFailedToAddUserToDatabase(err string) FailedToAddUserToDatabase {
	return FailedToAddUserToDatabase{
		msg: fmt.Sprintf("failed to add user to the database: %s", err),
	}
}

func (err FailedToAddUserToDatabase) Error() string {
	return err.msg
}

type FailedToGenerateToken struct {
	msg string
}

func CreateFailedToGenerateToken(err string) FailedToGenerateToken {
	return FailedToGenerateToken{
		msg: fmt.Sprintf("failed to generate new token: %s", err),
	}
}

func (err FailedToGenerateToken) Error() string {
	return err.msg
}

type FailedToCreateAuthenticationToken struct {
	msg string
}

func CreateFailedToCreateAuthenticationToken(err string) FailedToCreateAuthenticationToken {
	return FailedToCreateAuthenticationToken{
		msg: fmt.Sprintf("failed to create authentication token: %s", err),
	}
}

func (err FailedToCreateAuthenticationToken) Error() string {
	return err.msg
}
