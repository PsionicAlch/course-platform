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

	// Tutorial based functions.
	AddKeyword(keyword string) (string, error)
	GetAllTutorials() ([]*models.TutorialModel, error)
	AddNewTutorial(title, slug, description, thumbnailUrl, bannerUrl, content, fileChecksum string) error
	AddNewTutorialBulk(tutorials []*models.TutorialModel) error
	UpdateTutorial(id string, tutorial *models.TutorialModel) error
	UpdateTutorialBulk(tutorials []*models.TutorialModel) error
}
