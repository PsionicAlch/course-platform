package authentication

import (
	"context"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

type ContextKey string

const UserContextKey ContextKey = "user-context-key"

type RequestWithUser struct {
	*http.Request
	User *models.UserModel
}

func NewContextWithUserModel(user *models.UserModel, reqCtx context.Context) context.Context {
	return context.WithValue(reqCtx, UserContextKey, user)
}

func GetUserFromRequest(r *http.Request) *models.UserModel {
	user, ok := r.Context().Value(UserContextKey).(*models.UserModel)
	if !ok {
		return nil
	}

	return user
}
