package purchases

import (
	"net/http"
	"net/url"
	"slices"
	"strconv"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
	"github.com/justinas/nosurf"
)

const PurchasesPerPagination = 25

var PurchaseStatuses = []string{
	database.Pending.String(),
	database.RequiresAction.String(),
	database.Processing.String(),
	database.Succeeded.String(),
	database.Failed.String(),
	database.Cancelled.String(),
	database.Refunded.String(),
	database.Disputed.String(),
}

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
		Renderers: *pages.CreateRenderers(pageRenderer, htmxRenderer, nil),
		Database:  db,
		Auth:      auth,
	}
}

func (h *Handlers) PurchasesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminPurchasesPage{
		BasePage:        html.NewBasePage(user, nosurf.Token(r)),
		PaymentStatuses: PurchaseStatuses,
	}

	numPurchases, err := h.Database.CountAllPurchases()
	if err != nil {
		h.ErrorLog.Printf("Failed to count all course purchases: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.NumPurchases = numPurchases

	courses, err := h.Database.GetAllCourses("", nil)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all courses: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Courses = courses

	authors, err := h.Database.GetUsers("", database.Author, "", "")
	if err != nil {
		h.ErrorLog.Printf("Failed to get all authors: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Authors = authors

	users, err := h.Database.GetAllUsers()
	if err != nil {
		h.ErrorLog.Printf("Failed to get all users: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Users = users

	discounts, err := h.Database.GetAllDiscounts()
	if err != nil {
		h.ErrorLog.Printf("Failed to get all discounts: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Discounts = discounts

	coursePurchaseList, urlQuery, err := h.CreateCoursePurchasesList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create course list: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Purchases = coursePurchaseList

	urlQuery.Set("page", "1")
	pageData.URLQuery = urlQuery.Encode()

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-purchases", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PurchasesPaginationGet(w http.ResponseWriter, r *http.Request) {
	coursePurchaseList, _, err := h.CreateCoursePurchasesList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create course list: %s\n", err)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "admin-purchases", html.AdminCoursePurchaseListComponent{ErrorMessage: "Failed to get course purchases. Please try again"}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "admin-purchases", coursePurchaseList); err != nil {
		h.ErrorLog.Println(err)
	}
}

// Possible URL queries:
// -page
// -query
// -course
// -author
// -status
func (h *Handlers) CreateCoursePurchasesList(r *http.Request) (*html.AdminCoursePurchaseListComponent, url.Values, error) {
	var page int
	var query string
	var course string
	var author string
	var status string

	urlQuery := make(url.Values)

	if pageNum, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && page > 0 {
		page = pageNum
	} else {
		page = 1
	}

	urlQuery.Add("page", strconv.Itoa(page+1))

	if r.URL.Query().Get("query") != "" {
		query = r.URL.Query().Get("query")
		urlQuery.Add("query", query)
	}

	if r.URL.Query().Get("course") != "" {
		course = r.URL.Query().Get("course")
		urlQuery.Add("course", course)
	}

	if r.URL.Query().Get("author") != "" {
		author = r.URL.Query().Get("author")
		urlQuery.Add("course", author)
	}

	if s := r.URL.Query().Get("status"); s != "" && slices.Contains(PurchaseStatuses, s) {
		status = s
		urlQuery.Add("status", s)
	}

	coursePurchases, err := h.Database.AdminGetCoursePurchases(query, course, author, status, uint(page), PurchasesPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get course purchases (page %d): %s\n", page, err)
		return nil, urlQuery, err
	}

	var coursePurchasesSlice []*models.CoursePurchaseModel
	var lastCoursePurchase *models.CoursePurchaseModel

	if len(coursePurchases) < PurchasesPerPagination {
		coursePurchasesSlice = coursePurchases
	} else {
		coursePurchasesSlice = coursePurchases[:len(coursePurchases)-1]
		lastCoursePurchase = coursePurchases[len(coursePurchases)-1]
	}

	users := make(map[string]*models.UserModel)
	courses := make(map[string]*models.CourseModel)

	for _, purchase := range coursePurchases {
		if _, has := users[purchase.UserID]; !has {
			user, err := h.Database.GetUserByID(purchase.UserID, database.All)
			if err != nil {
				h.ErrorLog.Printf("Failed to get user by ID (\"%s\"): %s\n", purchase.UserID, err)
				return nil, urlQuery, err
			}

			users[purchase.UserID] = user
		}

		if _, has := courses[purchase.CourseID]; !has {
			course, err := h.Database.GetCourseByID(purchase.CourseID)
			if err != nil {
				h.ErrorLog.Printf("Failed to get course by ID (\"%s\"): %s\n", purchase.CourseID, err)
				return nil, urlQuery, err
			}

			courses[purchase.CourseID] = course
		}
	}

	coursePurchaseList := &html.AdminCoursePurchaseListComponent{
		Purchases:    coursePurchasesSlice,
		LastPurchase: lastCoursePurchase,
		Users:        users,
		Courses:      courses,
		BaseURL:      "/admin/purchases/htmx",
		URLQuery:     urlQuery.Encode(),
	}

	return coursePurchaseList, urlQuery, nil
}
