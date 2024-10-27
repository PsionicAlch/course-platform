package database

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/pkg/gatekeeper"
)

type Database interface {
	// General database functions.
	Close() error

	// Migration functions.
	MigrateUp() error
	MigrateDown() error
	Rollback(steps int) error

	// Authentication functions.
	GetUserInformation(email string) (*gatekeeper.UserInformation, error)
	AddUser(email, password string) (string, error)
	AddToken(token *gatekeeper.Token) error
	GetToken(token string) (*gatekeeper.Token, error)

	// User based functions.
	FindUserByEmail(email string) (*models.UserModel, error)
	FindUserByID(id string) (*models.UserModel, error)
}
