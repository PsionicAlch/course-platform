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

// ------------------------------------------
// -        Failed To Generate Key          -
// ------------------------------------------

type FailedToGenerateSecureCookieKey struct {
	msg string
}

func createFailedToGenerateSecureCookieKey(key, err string) FailedToGenerateSecureCookieKey {
	return FailedToGenerateSecureCookieKey{
		msg: fmt.Sprintf("failed to generate %s key: %s", key, err),
	}
}

func (err FailedToGenerateSecureCookieKey) Error() string {
	return err.msg
}

// --------------------------------------------------
// -        Failed To Encode Secure Cookie          -
// --------------------------------------------------

type FailedToEncodeSecureCookie struct {
	msg string
}

func createFailedToEncodeSecureCookie(err string) FailedToEncodeSecureCookie {
	return FailedToEncodeSecureCookie{
		msg: fmt.Sprintf("failed to encode secure cookie: %s", err),
	}
}

func (err FailedToEncodeSecureCookie) Error() string {
	return err.msg
}

// ------------------------------------------
// -        Invalid Gatekeeper Key          -
// ------------------------------------------

type InvalidGatekeeperKey struct {
	msg string
}

func createInvalidGatekeeperKey(err string) InvalidGatekeeperKey {
	var msg string
	if err == "" {
		msg = "invalid gatekeeper key was provided. Consider generating your key using gatekeeper.GenerateGatekeeperKey"
	} else {
		msg = fmt.Sprintf("invalid gatekeeper key was provided: %s\n. Consider generating your key using gatekeeper.GenerateGatekeeperKey", err)
	}

	return InvalidGatekeeperKey{
		msg: msg,
	}
}

func (err InvalidGatekeeperKey) Error() string {
	return err.msg
}
