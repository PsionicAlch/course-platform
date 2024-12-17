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
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
	"github.com/go-chi/chi/v5"
)

const TutorialsPerPagination = 25

type Handlers struct {
	utils.Loggers
	Renderers pages.Renderers
	Database  database.Database
}

func SetupHandlers(pageRenderer, htmxRenderer render.Renderer, db database.Database) *Handlers {
	loggers := utils.CreateLoggers("AUTHOR TUTORIAL HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: *pages.CreateRenderers(pageRenderer, htmxRenderer, nil),
		Database:  db,
	}
}

func (h *Handlers) TutorialsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AuthorsTutorialsPage{
		BasePage: html.NewBasePage(user),
	}

	tutorialsList, author, err := h.CreateTutorialsList(r)
	if err != nil {
		if err == ErrAuthorNotFound || err == ErrTutorialsNotFound {
			if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user)}, http.StatusNotFound); err != nil {
				h.ErrorLog.Println(err)
			}
		} else {
			h.ErrorLog.Printf("Failed to create tutorials list: %s\n", err)

			if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}
		}

		return
	}

	pageData.Author = author
	pageData.Tutorials = tutorialsList

	lenTutorials, err := h.Database.CountTutorialsWrittenBy(author.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to count all tutorials written by \"%s\": %s\n", author.ID, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.LenTutorials = lenTutorials

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "authors-tutorials", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialsPaginationGet(w http.ResponseWriter, r *http.Request) {
	tutorialsList, _, err := h.CreateTutorialsList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create tutorials list: %s\n", err)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "tutorials", html.TutorialsListComponent{ErrorMessage: "Failed to get tutorials. Please try again."}); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "tutorials", tutorialsList); err != nil {
		h.ErrorLog.Println(err)
	}
}

// Possible URL queries:
// -page
// -query
func (h *Handlers) CreateTutorialsList(r *http.Request) (*html.TutorialsListComponent, *models.UserModel, error) {
	authorSlug := chi.URLParam(r, "author-slug")

	query := r.URL.Query().Get("query")
	page := 1

	urlQuery := make(url.Values)

	if pageNum, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		page = pageNum
	}

	urlQuery.Add("query", query)
	urlQuery.Add("page", fmt.Sprintf("%d", page+1))

	author, err := h.Database.GetUserBySlug(authorSlug, database.Author)
	if err != nil {
		h.ErrorLog.Printf("Failed to find author by slug (\"%s\"): %s\n", authorSlug, err)
		return nil, nil, err
	}

	if author == nil {
		return nil, nil, ErrAuthorNotFound
	}

	tutorials, err := h.Database.GetTutorials(query, author.ID, page, TutorialsPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all tutorials (page %d): %s\n", page, err)
		return nil, nil, err
	}

	if len(tutorials) == 0 {
		return nil, nil, ErrTutorialsNotFound
	}

	var tutorialsSlice []*models.TutorialModel
	var lastTutorial *models.TutorialModel

	if len(tutorials) < TutorialsPerPagination {
		tutorialsSlice = tutorials
	} else {
		tutorialsSlice = tutorials[:len(tutorials)-1]
		lastTutorial = tutorials[len(tutorials)-1]
	}

	tutorialList := &html.TutorialsListComponent{
		Tutorials:    tutorialsSlice,
		LastTutorial: lastTutorial,
		QueryURL:     fmt.Sprintf("/authors/%s/tutorials/htmx?%s", author.Slug, urlQuery.Encode()),
	}

	return tutorialList, author, nil
}
