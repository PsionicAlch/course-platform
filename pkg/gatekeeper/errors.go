package gatekeeper

import "fmt"

// ----------------------------------------------
// -		Failed To Find User By Email		-
// ----------------------------------------------

type FailedToFindUserByEmail struct {
	msg string
}

func createFailedToFindUserByEmail(email, err string) FailedToFindUserByEmail {
	return FailedToFindUserByEmail{
		msg: fmt.Sprintf("failed to find user by %s: %s", email, err),
	}
}

func (err FailedToFindUserByEmail) Error() string {
	return err.msg
}

// --------------------------------------
// -		User Already Exists			-
// --------------------------------------

type UserAlreadyExists struct {
	msg string
}

func createUserAlreadyExists(email string) UserAlreadyExists {
	return UserAlreadyExists{
		msg: fmt.Sprintf("a user with the email %s already exists", email),
	}
}

func (err UserAlreadyExists) Error() string {
	return err.msg
}

// ------------------------------------------
// -		Failed To Hash Password			-
// ------------------------------------------

type FailedToHashPassword struct {
	msg string
}

func createFailedToHashPassword(err string) FailedToHashPassword {
	return FailedToHashPassword{
		msg: fmt.Sprintf("failed to hash user's password: %s", err),
	}
}

func (err FailedToHashPassword) Error() string {
	return err.msg
}

// --------------------------------------------------
// -		Failed To Add User To Database			-
// --------------------------------------------------

type FailedToAddUserToDatabase struct {
	msg string
}

func createFailedToAddUserToDatabase(err string) FailedToAddUserToDatabase {
	return FailedToAddUserToDatabase{
		msg: fmt.Sprintf("failed to add user to the database: %s", err),
	}
}

func (err FailedToAddUserToDatabase) Error() string {
	return err.msg
}

// ----------------------------------------------
// -		Failed To Generate New Token		-
// ----------------------------------------------

type FailedToGenerateNewToken struct {
	msg string
}

func createFailedToGenerateNewToken(err string) FailedToGenerateNewToken {
	return FailedToGenerateNewToken{
		msg: fmt.Sprintf("failed to generate new token: %s", err),
	}
}

func (err FailedToGenerateNewToken) Error() string {
	return err.msg
}

// --------------------------------------------------
// -		Failed To Add Token To Database			-
// --------------------------------------------------

type FailedToAddTokenToDatabase struct {
	msg string
}

func createFailedToAddTokenToDatabase(err string) FailedToAddTokenToDatabase {
	return FailedToAddTokenToDatabase{
		msg: fmt.Sprintf("failed to add token to the database: %s", err),
	}
}

func (err FailedToAddTokenToDatabase) Error() string {
	return err.msg
}

// ----------------------------------------------------------
// -		Failed To Create Authentication Cookie			-
// ----------------------------------------------------------

type FailedToCreateAuthenticationCookie struct {
	msg string
}

func createFailedToCreateAuthenticationCookie(err string) FailedToCreateAuthenticationCookie {
	return FailedToCreateAuthenticationCookie{
		msg: fmt.Sprintf("failed to create an authentication cookie: %s", err),
	}
}

func (err FailedToCreateAuthenticationCookie) Error() string {
	return err.msg
}
