package database

import "time"

type Database interface {
	// General database functions.
	Close() error

	// Migration functions.
	MigrateUp() error
	MigrateDown() error
	Rollback(steps int) error

	// Users functions.
	UserExists(email string) (bool, error)
	AddUser(name, surname, email, password string) (string, error)

	// Tokens functions.
	AddToken(token, tokenType, userId, ipAddr string, validUntil time.Time) error
}
