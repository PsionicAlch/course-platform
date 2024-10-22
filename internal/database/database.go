package database

type Database interface {
	// General database functions.
	Close() error
	MigrateUp() error
	MigrateDown() error
	Rollback(steps int) error

	// Authentication functions.
	AddUser(email, password string) (string, error)
	CreateAuthenticationToken(userId, ip string) (string, error)
}
