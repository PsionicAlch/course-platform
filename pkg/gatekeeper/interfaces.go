package gatekeeper

type GatekeeperDatabase interface {
	// UserExists takes a user string and returns a string that represents the
	// user's ID and possibly and error if anything went wrong. This function's
	// responsibility is to query the database for a user using the user's email address
	// and then return the ID of the user (in string format to support all possible ID
	// types) or an empty string if that user couldn't be found. The function can also return
	// an error if there were any.
	UserExists(email string) (bool, error)

	// AddUser takes in the user's email address and the user's hashed password and returns
	// the user's ID as well as a possible error. This function's responsibility is to add
	// a new user to the database and return their ID (in string format to support all
	// possible ID types).
	AddUser(email, password string) (string, error)

	// AddToken takes in a token, a token type (to distinguish auth tokens form email tokens),
	// validUntil (the date and time at which this token expires), the user's ID to which
	// this token is connected, and the IP address from which this user logged in. This
	// function's responsibility is to add the new token, and associated data, to your
	// database for later retrieval.
	AddToken(token *Token) error

	// TokenExists takes in a token and it's type then returns a bool or an error. This
	// function's responsibility is to check if a token with that type exists in the database
	GetToken(token string) (*Token, error)
}
