package tutorials

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
	"github.com/PsionicAlch/psionicalch-home/web/html"
	"github.com/PsionicAlch/psionicalch-home/web/pages"
	"github.com/justinas/nosurf"
)

const TutorialsPerPagination = 25

type Handlers struct {
	utils.Loggers
	Renderers pages.Renderers
	Database  database.Database
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, db database.Database) *Handlers {
	loggers := utils.CreateLoggers("PROFILE HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: *pages.CreateRenderers(pageRenderer, htmxRenderer, nil),
		Database:  db,
	}
}

func (h *Handlers) TutorialsBookmarksGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.ProfileTutorialsBookmarksPage{
		BasePage: html.NewBasePage(user, nosurf.Token(r)),
	}

	tutorialsList, err := h.CreateTutorialsList(r, "/profile/tutorials/bookmarks/htmx", h.Database.GetTutorialsBookmarkedByUser)
	if err != nil {
		h.ErrorLog.Printf("Failed to create tutorials list: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Tutorials = tutorialsList

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "profile-tutorials-bookmarks", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialsBookmarksPaginationGet(w http.ResponseWriter, r *http.Request) {
	tutorialsList, err := h.CreateTutorialsList(r, "/profile/tutorials/bookmarks/htmx", h.Database.GetTutorialsBookmarkedByUser)
	if err != nil {
		h.ErrorLog.Printf("Failed to create tutorials list: %s\n", err)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "tutorials", html.TutorialsListComponent{ErrorMessage: "Failed to get tutorials. Please try again."}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "tutorials", tutorialsList); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialsLikesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.ProfileTutorialsLikedPage{
		BasePage: html.NewBasePage(user, nosurf.Token(r)),
	}

	tutorialsList, err := h.CreateTutorialsList(r, "/profile/tutorials/likes/htmx", h.Database.GetTutorialsLikedByUser)
	if err != nil {
		h.ErrorLog.Printf("Failed to create tutorials list: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Tutorials = tutorialsList

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "profile-tutorials-likes", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialsLikesPaginationGet(w http.ResponseWriter, r *http.Request) {
	tutorialsList, err := h.CreateTutorialsList(r, "/profile/tutorials/likes/htmx", h.Database.GetTutorialsLikedByUser)
	if err != nil {
		h.ErrorLog.Printf("Failed to create tutorials list: %s\n", err)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "tutorials", html.TutorialsListComponent{ErrorMessage: "Failed to get tutorials. Please try again."}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "tutorials", tutorialsList); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CreateTutorialsList(r *http.Request, baseURL string, dbFunc func(term string, userId string, page uint, elements uint) ([]*models.TutorialModel, error)) (*html.TutorialsListComponent, error) {
	user := authentication.GetUserFromRequest(r)

	page := 1
	query := ""

	urlQuery := make(url.Values)

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		page = p
	}

	urlQuery.Add("page", fmt.Sprintf("%d", page+1))

	if q := r.URL.Query().Get("query"); q != "" {
		query = q

		urlQuery.Add("query", q)
	}

	tutorials, err := dbFunc(query, user.ID, uint(page), TutorialsPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all tutorials liked/bookmarked by user (\"%s\"): %s\n", user.ID, err)
		return nil, err
	}

	var tutorialsSlice []*models.TutorialModel
	var lastTutorial *models.TutorialModel

	if len(tutorials) < TutorialsPerPagination {
		tutorialsSlice = tutorials
	} else {
		tutorialsSlice = tutorials[:len(tutorials)-1]
		lastTutorial = tutorials[len(tutorials)-1]
	}

	tutorialsList := &html.TutorialsListComponent{
		Tutorials:    tutorialsSlice,
		LastTutorial: lastTutorial,
		QueryURL:     fmt.Sprintf("%s?%s", baseURL, urlQuery.Encode()),
	}

	return tutorialsList, nil
}
