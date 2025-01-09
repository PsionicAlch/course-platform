package refunds

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
	"github.com/PsionicAlch/psionicalch-home/web/html"
	"github.com/PsionicAlch/psionicalch-home/web/pages"
	"github.com/justinas/nosurf"
)

const RefundsPerPagination = 25

var RefundStatuses = []string{
	database.RefundPending.String(),
	database.RefundRequiresAction.String(),
	database.RefundSucceeded.String(),
	database.RefundCancelled.String(),
	database.DisputeWarningNeedsResponse.String(),
	database.DisputeWarningUnderReview.String(),
	database.DisputeWarningClosed.String(),
	database.DisputeNeedsResponse.String(),
	database.DisputeUnderReview.String(),
	database.DisputeWon.String(),
	database.DisputeLost.String(),
}

type Handlers struct {
	utils.Loggers
	Render   pages.Renderers
	Database database.Database
	Auth     *authentication.Authentication
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, db database.Database, auth *authentication.Authentication) *Handlers {
	loggers := utils.CreateLoggers("ADMIN HANDLERS")

	return &Handlers{
		Loggers:  loggers,
		Render:   *pages.CreateRenderers(pageRenderer, htmxRenderer, nil),
		Database: db,
		Auth:     auth,
	}
}

func (h *Handlers) RefundsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminRefundsPage{
		BasePage:       html.NewBasePage(user, nosurf.Token(r)),
		RefundStatuses: RefundStatuses,
	}

	numRefunds, err := h.Database.CountRefunds()
	if err != nil {
		h.ErrorLog.Printf("Failed to count the number of refunds in the database: %s\n", err)

		if err := h.Render.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.NumRefunds = numRefunds

	refunds, urlQuery, err := h.CreateRefundsList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create refunds list: %s\n", err)

		if err := h.Render.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Refunds = refunds

	urlQuery.Set("page", "1")
	pageData.URLQuery = urlQuery.Encode()

	if err := h.Render.Page.RenderHTML(w, r.Context(), "admin-refunds", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) RefundsPaginationGet(w http.ResponseWriter, r *http.Request) {
	refundsList, _, err := h.CreateRefundsList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create refunds list: %s\n", err)

		if err := h.Render.Htmx.RenderHTML(w, nil, "admin-refunds", html.AdminRefundsListComponent{ErrorMessage: "Failed to get refunds. Please try again."}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Render.Htmx.RenderHTML(w, nil, "admin-refunds", refundsList); err != nil {
		h.ErrorLog.Println(err)
	}
}

// Possible URL Queries:
// - page
// - query
// - status
func (h *Handlers) CreateRefundsList(r *http.Request) (*html.AdminRefundsListComponent, url.Values, error) {
	var page uint
	var query string
	var status string

	urlQuery := make(url.Values)

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		page = uint(p)
	} else {
		page = 1
	}

	urlQuery.Add("page", strconv.Itoa(int(page+1)))

	if q := r.URL.Query().Get("query"); q != "" {
		query = q
		urlQuery.Add("query", q)
	}

	if s := r.URL.Query().Get("status"); slices.Contains(RefundStatuses, s) {
		status = s
		urlQuery.Add("status", s)
	}

	refunds, err := h.Database.AdminGetRefunds(query, status, uint(page), RefundsPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all refunds from the database: %s\n", err)
		return nil, urlQuery, err
	}

	var refundsSlice []*models.RefundModel
	var lastRefund *models.RefundModel

	if len(refunds) < RefundsPerPagination {
		refundsSlice = refunds
	} else {
		refundsSlice = refunds[:len(refunds)-1]
		lastRefund = refunds[len(refunds)-1]
	}

	users := make(map[string]*models.UserModel)
	courses := make(map[string]*models.CourseModel)

	for _, refund := range refunds {
		if _, has := users[refund.UserID]; !has {
			user, err := h.Database.GetUserByID(refund.UserID, database.All)
			if err != nil {
				h.ErrorLog.Printf("Failed to get user by ID (\"%s\"): %s\n", refund.UserID, err)
				return nil, urlQuery, err
			}

			users[refund.UserID] = user
		}

		if _, has := courses[refund.CoursePurchaseID]; !has {
			course, err := h.Database.GetCourseByCoursePurchaseID(refund.CoursePurchaseID)
			if err != nil {
				h.ErrorLog.Printf("Failed to get course by course purchase ID (\"%s\"): %s\n", refund.CoursePurchaseID, err)
				return nil, urlQuery, err
			}

			courses[refund.CoursePurchaseID] = course
		}
	}

	refundsList := &html.AdminRefundsListComponent{
		Refunds:    refundsSlice,
		LastRefund: lastRefund,
		Users:      users,
		Courses:    courses,
		BaseURL:    "/admin/refunds/htmx",
		URLQuery:   urlQuery.Encode(),
	}

	return refundsList, urlQuery, nil
}
