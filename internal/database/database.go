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
	AddUser(name, surname, email, password string) (*models.UserModel, error)
	GetUser(email string) (*models.UserModel, error)
	GetUserByID(id string) (*models.UserModel, error)
	GetUserByToken(token, tokenType string) (*models.UserModel, error)
	UpdateUserPassword(userId, password string) error

	// Tokens functions.
	AddToken(token, tokenType, userId, ipAddr string, validUntil time.Time) error
	GetToken(token, tokenType string) (*models.TokenModel, error)
	DeleteToken(token, tokenType string) error
}
