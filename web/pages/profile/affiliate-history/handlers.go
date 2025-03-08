package affiliatehistory

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PsionicAlch/course-platform/internal/authentication"
	"github.com/PsionicAlch/course-platform/internal/database/models"
	"github.com/PsionicAlch/course-platform/internal/utils"
	"github.com/PsionicAlch/course-platform/web/html"
	"github.com/PsionicAlch/course-platform/web/pages"
	"github.com/justinas/nosurf"
)

const AffiliateHistoryElementsPerPagination = 25

type Handlers struct {
	utils.Loggers
	*pages.HandlerContext
}

func SetupHandlers(handlerContext *pages.HandlerContext) *Handlers {
	loggers := utils.CreateLoggers("PROFILE AFFILIATE HISTORY HANDLERS")

	return &Handlers{
		Loggers:        loggers,
		HandlerContext: handlerContext,
	}
}

func (h *Handlers) AffiliateHistoryGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.ProfileAffiliateHistoryPage{
		BasePage: html.NewBasePage(user, nosurf.Token(r)),
		User:     user,
	}

	affiliateHistory, err := h.CreateAffiliateHistoryList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create affiliate history list: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.AffiliateHistory = affiliateHistory

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "profile-affiliate-history", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) AffiliateHistoryPaginationGet(w http.ResponseWriter, r *http.Request) {
	affiliateHistory, err := h.CreateAffiliateHistoryList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create affiliate history list: %s\n", err)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "affiliate-history", html.AffiliateHistoryListComponent{ErrorMessage: "Failed to get affiliate history. Please try again."}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "affiliate-history", affiliateHistory); err != nil {
		h.ErrorLog.Println(err)
	}
}

// Possible URL Query:
// - page
func (h *Handlers) CreateAffiliateHistoryList(r *http.Request) (*html.AffiliateHistoryListComponent, error) {
	user := authentication.GetUserFromRequest(r)
	page := 1

	urlQuery := make(url.Values)

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		page = p
		urlQuery.Add("page", fmt.Sprintf("%d", page+1))
	}

	affiliateHistory, err := h.Database.GetUserAffiliatePointsHistory(user.ID, uint(page), AffiliateHistoryElementsPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get user's (\"%s\") affiliate points history: %s\n", user.ID, err)
		return nil, err
	}

	var affiliateHistorySlice []*models.AffiliatePointsHistoryModel
	var lastAffiliateHistory *models.AffiliatePointsHistoryModel

	if len(affiliateHistory) < AffiliateHistoryElementsPerPagination {
		affiliateHistorySlice = affiliateHistory
	} else {
		affiliateHistorySlice = affiliateHistory[:len(affiliateHistory)-1]
		lastAffiliateHistory = affiliateHistory[len(affiliateHistory)-1]
	}

	affiliateHistoryList := &html.AffiliateHistoryListComponent{
		AffiliateHistory:     affiliateHistorySlice,
		LastAffiliateHistory: lastAffiliateHistory,
		QueryURL:             fmt.Sprintf("/profile/affiliate-history/htmx?%s", urlQuery.Encode()),
	}

	return affiliateHistoryList, nil
}
