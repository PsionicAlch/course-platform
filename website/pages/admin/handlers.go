package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	Renderers pages.Renderers
	Database  database.Database
	Auth      *authentication.Authentication
}

const UsersPerPagination = 5

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, db database.Database, auth *authentication.Authentication) *Handlers {
	loggers := utils.CreateLoggers("ADMIN HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: *pages.CreateRenderers(pageRenderer, htmxRenderer),
		Database:  db,
		Auth:      auth,
	}
}

func (h *Handlers) AdminGet(w http.ResponseWriter, r *http.Request) {
	utils.Redirect(w, r, "/admin/admins")
}

func (h *Handlers) AdminsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminAdminsPage{
		BasePage: html.NewBasePage(user),
	}

	adminUsers, err := h.Database.GetAllAdminsPaginated(1, UsersPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get admin users (page %d) from the database: %s\n", 1, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user),
		}); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var adminUsersSlice []*models.UserModel
	var lastAdminUser *models.UserModel

	if len(adminUsers) < UsersPerPagination {
		adminUsersSlice = adminUsers
	} else {
		adminUsersSlice = adminUsers[:len(adminUsers)-1]
		lastAdminUser = adminUsers[len(adminUsers)-1]
	}

	pageData.Admins = &html.AdminUsersListComponent{
		Admins:    adminUsersSlice,
		LastAdmin: lastAdminUser,
		QueryURL:  fmt.Sprintf("/admin/admins/htmx?page=%d", 2),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-admins", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) AdminsPaginationGet(w http.ResponseWriter, r *http.Request) {
	adminsList := html.AdminUsersListComponent{}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		h.ErrorLog.Printf("Failed to convert page to int: %s\n", err)

		adminsList.ErrorMessage = "Failed to get admin users. Please try again."
		if err := h.Renderers.Htmx.RenderHTML(w, nil, "admins", adminsList, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	adminUsers, err := h.Database.GetAllAdminsPaginated(page, UsersPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get admin users (page %d) from the database: %s\n", page, err)

		adminsList.ErrorMessage = "Failed to get admin users. Please try again."
		if err := h.Renderers.Htmx.RenderHTML(w, nil, "admins", adminsList, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var adminUsersSlice []*models.UserModel
	var lastAdminUser *models.UserModel

	if len(adminUsers) < UsersPerPagination {
		adminUsersSlice = adminUsers
	} else {
		adminUsersSlice = adminUsers[:len(adminUsers)-1]
		lastAdminUser = adminUsers[len(adminUsers)-1]
	}

	adminsList.Admins = adminUsersSlice
	adminsList.LastAdmin = lastAdminUser
	adminsList.QueryURL = fmt.Sprintf("/admin/admins/htmx?page=%d", page+1)

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "admins", adminsList); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) AuthorsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminAuthorsPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-authors", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CommentsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminCommentsPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-comments", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CoursesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminCoursesPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-courses", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

// TODO: Consider adding a usage_amount to discounts so that they can only be used a set amount of times.

func (h *Handlers) DiscountsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminDiscountsPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-discounts", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PurchasesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminPurchasesPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-purchases", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) RefundsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminRefundsPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-refunds", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminTutorialsPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-tutorials", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) UsersGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminUsersPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-users", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}
