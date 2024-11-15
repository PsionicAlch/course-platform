package admin

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	renderers pages.Renderers
	auth      *authentication.Authentication
}

func SetupHandlers(pageRenderer render.Renderer, auth *authentication.Authentication) *Handlers {
	loggers := utils.CreateLoggers("ADMIN HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		renderers: *pages.CreateRenderers(pageRenderer, nil),
		auth:      auth,
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

	if err := h.renderers.Page.RenderHTML(w, "admin-admins", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) AuthorsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminAuthorsPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, "admin-authors", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CommentsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminCommentsPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, "admin-comments", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CoursesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminCoursesPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, "admin-courses", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

// TODO: Consider adding a usage_amount to discounts so that they can only be used a set amount of times.

func (h *Handlers) DiscountsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminDiscountsPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, "admin-discounts", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PurchasesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminPurchasesPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, "admin-purchases", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) RefundsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminRefundsPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, "admin-refunds", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminTutorialsPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, "admin-tutorials", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) UsersGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminUsersPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.renderers.Page.RenderHTML(w, "admin-users", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}
