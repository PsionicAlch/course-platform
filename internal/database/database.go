package database

type Database interface {
	// General database functions.
	Close()

	// Authentication functions.
	AddUser(email, password string) (string, error)
	CreateAuthenticationToken(userId string) (string, error)
}
