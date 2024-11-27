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
	"github.com/go-chi/chi/v5"
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

	usersCount, err := h.Database.CountUsers()
	if err != nil {
		h.ErrorLog.Printf("Failed to get the number of users from the database: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.NumUsers = usersCount

	pageData.AuthorizationLevels = database.AuthorizationLevelStrings()

	usersList, err := h.CreateUsersList("", database.All, 1)
	if err != nil {
		h.ErrorLog.Printf("Failed to get the number of users from the database: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Users = usersList

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-users", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) UsersPaginationGet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	level := database.All
	page := 2

	if pageNum, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		page = pageNum
	}

	if authLevel, err := database.AuthorizationLevelString(r.URL.Query().Get("level")); err == nil {
		level = authLevel
	}

	usersList, err := h.CreateUsersList(query, level, uint(page))
	if err != nil {
		if err := h.Renderers.Htmx.RenderHTML(w, nil, "users", &html.UsersListComponent{
			ErrorMessage: "Failed to get users. Please try again.",
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "users", usersList); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) AuthorEditGet(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "user-id")

	user, err := h.Database.GetUserByID(userId)
	if err != nil {
		h.ErrorLog.Printf("Failed to get user by ID \"%s\": %s\n", userId, err)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", html.SelectComponent{
			Name: "author-status",
			Options: []string{
				"Author",
				"User",
			},
			ErrorMessage: "Unable to find user by given ID.",
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var userStatus string
	if user.IsAuthor {
		userStatus = "Author"
	} else {
		userStatus = "User"
	}

	selectComponent := html.SelectComponent{
		Name: "author-status",
		Options: []string{
			"Author",
			"User",
		},
		Selected: userStatus,
		URL:      fmt.Sprintf("/admin/users/htmx/change-author/%s", user.ID),
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", selectComponent); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) AuthorEditPost(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "user-id")

	user, err := h.Database.GetUserByID(userId)
	if err != nil {
		h.ErrorLog.Printf("Failed to get user by ID \"%s\": %s\n", userId, err)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", html.SelectComponent{
			Name: "author-status",
			Options: []string{
				"Author",
				"User",
			},
			ErrorMessage: "Unable to find user by given ID.",
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	r.ParseForm()

	authorStatus := r.Form.Get("author-status")

	if !utils.InSlice(authorStatus, []string{"Author", "User"}) {
		h.ErrorLog.Printf("Received invalid author status option: %s\n", authorStatus)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", html.SelectComponent{
			Name: "author-status",
			Options: []string{
				"Author",
				"User",
			},
			ErrorMessage: "Invalid option selected.",
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if authorStatus == "Author" {
		if err := h.Database.AddAuthorStatus(user.ID); err != nil {
			h.ErrorLog.Printf("Failed to update user's (\"%s\") author status: %s\n", user.ID, err)

			if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", html.SelectComponent{
				Name: "author-status",
				Options: []string{
					"Author",
					"User",
				},
				ErrorMessage: "Unexpected server error has occurred!",
			}, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", "<p style=\"color: var(--primary-green-color);\">&#10004;</p>"); err != nil {
			h.ErrorLog.Println(err)
		}
	} else {
		if err := h.Database.RemoveAuthorStatus(user.ID); err != nil {
			h.ErrorLog.Printf("Failed to update user's (\"%s\") author status: %s\n", user.ID, err)

			if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", html.SelectComponent{
				Name: "author-status",
				Options: []string{
					"Author",
					"User",
				},
				ErrorMessage: "Unexpected server error has occurred!",
			}, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", "<p style=\"color: var(--primary-red-color);\">&#10008;</p>"); err != nil {
			h.ErrorLog.Println(err)
		}
	}
}

func (h *Handlers) CreateUsersList(term string, level database.AuthorizationLevel, page uint) (*html.UsersListComponent, error) {
	users, err := h.Database.GetUsersPaginated(term, level, page, UsersPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get users (page %d) from the database: %s\n", 1, err)

		return nil, err
	}

	tutorialsLiked := make(map[string]uint, len(users))
	tutorialsBookmarked := make(map[string]uint, len(users))
	coursesBought := make(map[string]uint, len(users))
	tutorialsWritten := make(map[string]uint, len(users))
	coursesWritten := make(map[string]uint, len(users))

	for _, user := range users {
		tutsLiked, err := h.Database.CountTutorialsLikedByUser(user.ID)
		if err != nil {
			h.ErrorLog.Printf("Failed to get the number of tutorials liked by user \"%s\": %s\n", user.ID, err)

			return nil, err
		}

		tutorialsLiked[user.ID] = tutsLiked

		tutorialsBookmarked[user.ID] = 0

		coursesBought[user.ID] = 0

		tutorialsWritten[user.ID] = 0

		coursesWritten[user.ID] = 0
	}

	var usersSlice []*models.UserModel
	var lastUser *models.UserModel

	if len(users) < UsersPerPagination {
		usersSlice = users
	} else {
		usersSlice = users[:len(users)-1]
		lastUser = users[len(users)-1]
	}

	usersList := &html.UsersListComponent{
		Users:               usersSlice,
		LastUser:            lastUser,
		TutorialsLiked:      tutorialsLiked,
		TutorialsBookmarked: tutorialsBookmarked,
		CoursesBought:       coursesBought,
		TutorialsWritten:    tutorialsWritten,
		CoursesWritten:      coursesWritten,
		BaseURL:             "/admin/users/htmx",
		QueryURL:            fmt.Sprintf("/admin/users/htmx?query=%s&level=%s&page=%d", url.QueryEscape(term), url.QueryEscape(level.String()), page+1),
	}

	return usersList, nil
}
