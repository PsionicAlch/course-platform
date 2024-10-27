package pages

import (
	"fmt"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/pkg/gatekeeper"
)

func GetUser(cookies []*http.Cookie, auth *gatekeeper.Gatekeeper, db database.Database) (*models.UserModel, error) {
	id, err := auth.GetUserIDFromAuthenticationToken(cookies)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id from authentication token: %s", err)
	}

	if id == "" {
		return nil, nil
	}

	user, err := db.FindUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from database: %s", err)
	}

	return user, nil
}
