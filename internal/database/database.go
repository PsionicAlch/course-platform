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

	// Users functions.
	UserExists(email string) (bool, error)
	AddUser(name, surname, email, password string) (string, error)
	GetUser(email string) (*models.UserModel, error)
	GetUserByID(id string) (*models.UserModel, error)

	// Tokens functions.
	AddToken(token, tokenType, userId, ipAddr string, validUntil time.Time) error
	GetToken(token, tokenType string) (*models.TokenModel, error)
}
