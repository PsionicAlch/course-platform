package database

import (
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

type Database interface {
	// General database functions.
	Close() error

	// Migration functions.
	MigrateUp() error
	MigrateDown() error
	Rollback(steps int) error

	// Authentication functions.
	UserExists(email string) (bool, error)
	AddUser(email, password string) (string, error)
	AddToken(token, tokenType string, validUntil time.Time, userId, ipAddr string) error

	// User based functions.
	FindUserByEmail(email string) (*models.UserModel, error)
}
