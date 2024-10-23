package database

import "github.com/PsionicAlch/psionicalch-home/internal/database/models"

type Database interface {
	// General database functions.
	Close() error

	// Migration functions.
	MigrateUp() error
	MigrateDown() error
	Rollback(steps int) error

	// Authentication functions.
	AddUser(email, password string) (string, error)
	CreateAuthenticationToken(userId, ip string) (string, error)

	// User based functions.
	FindUserByEmail(email string) (*models.UserModel, error)
}
