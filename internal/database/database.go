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
	AddNewUser(name, surname, email, password, token, tokenType, ipAddr string, validUntil time.Time) (*models.UserModel, error)
	GetUser(email string) (*models.UserModel, error)
	GetUserByID(id string) (*models.UserModel, error)
	GetUserByToken(token, tokenType string) (*models.UserModel, error)
	UpdateUserPassword(userId, password string) error

	// Tokens functions.
	AddToken(token, tokenType, userId string, validUntil time.Time) error
	GetToken(token, tokenType string) (*models.TokenModel, error)
	DeleteToken(token, tokenType string) error
	DeleteAllTokens(email, tokenType string) error

	// IP Addresses functions.
	AddIPAddress(userId, ipAddr string) error
	GetUserIpAddresses(userId string) ([]string, error)

	// Tutorials functions.
	GetAllTutorials() ([]*models.TutorialModel, error)
	GetAllTutorialsPaginated(page, elements int) ([]*models.TutorialModel, error)
	SearchTutorialsPaginated(term string, page, elements int) ([]*models.TutorialModel, error)
	GetTutorialBySlug(slug string) (*models.TutorialModel, error)
	BulkAddTutorials(tutorials []*models.TutorialModel) error
	BulkUpdateTutorials(tutorials []*models.TutorialModel) error
}
