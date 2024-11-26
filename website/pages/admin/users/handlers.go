package users

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

const UsersPerPagination = 25

type Handlers struct {
	utils.Loggers
	Renderers pages.Renderers
	Database  database.Database
	Auth      *authentication.Authentication
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, db database.Database, auth *authentication.Authentication) *Handlers {
	loggers := utils.CreateLoggers("ADMIN HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: *pages.CreateRenderers(pageRenderer, htmxRenderer),
		Database:  db,
		Auth:      auth,
	}
}

func (h *Handlers) UsersGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminUsersPage{
		BasePage: html.NewBasePage(user),
	}

	users, err := h.Database.GetUsersPaginated("", database.All, 1, UsersPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get users (page %d) from the database: %s\n", 1, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var usersSlice []*models.UserModel
	var lastUser *models.UserModel

	if len(users) < UsersPerPagination {
		usersSlice = users
	} else {
		usersSlice = users[:len(users)-1]
		lastUser = users[len(users)-1]
	}

	pageData.Users = &html.UsersListComponent{
		Users:    usersSlice,
		LastUser: lastUser,
		QueryURL: fmt.Sprintf("/admin/users/htmx?page=%d", 2),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-users", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) UsersPaginationGet(w http.ResponseWriter, r *http.Request) {
	usersList := &html.UsersListComponent{}

	query := r.URL.Query().Get("query")
	level := database.All
	page := 2

	if pageNum, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		page = pageNum
	}

	if authLevel, err := database.AuthorizationLevelString(r.URL.Query().Get("level")); err == nil {
		level = authLevel
	}

	users, err := h.Database.GetUsersPaginated(query, level, uint(page), UsersPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get users (page %d) from the database: %s\n", 1, err)
		usersList.ErrorMessage = "Failed to get users. Please try again."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "users", usersList, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var usersSlice []*models.UserModel
	var lastUser *models.UserModel

	if len(users) < UsersPerPagination {
		usersSlice = users
	} else {
		usersSlice = users[:len(users)-1]
		lastUser = users[len(users)-1]
	}

	usersList.Users = usersSlice
	usersList.LastUser = lastUser
	usersList.QueryURL = fmt.Sprintf("/admin/users/htmx?query=%s&level=%s&page=%d", url.QueryEscape(query), url.QueryEscape(level.String()), page+1)

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "users", usersList); err != nil {
		h.ErrorLog.Println(err)
	}
}
