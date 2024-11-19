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
	"github.com/PsionicAlch/psionicalch-home/internal/session"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
	"github.com/go-chi/chi/v5"
)

const TutorialsPerPagination = 5

type Handlers struct {
	utils.Loggers
	Renderers *pages.Renderers
	Database  database.Database
	Session   *session.Session
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, db database.Database, sessions *session.Session) *Handlers {
	loggers := utils.CreateLoggers("TUTORIALS HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: pages.CreateRenderers(pageRenderer, htmxRenderer),
		Database:  db,
		Session:   sessions,
	}
}

func (h *Handlers) TutorialsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.TutorialsPage{
		BasePage: html.NewBasePage(user),
	}

	tutorials, err := h.Database.GetAllTutorialsPaginated(1, TutorialsPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all tutorials (page 1) from the database: %s\n", err)

		h.Session.SetErrorMessage(r.Context(), "Failed to load tutorials. Please try again")

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "tutorials", pageData); err != nil {
			h.ErrorLog.Println(err)
		}
	}

	var tutSlice []*models.TutorialModel
	var lastTut *models.TutorialModel

	if len(tutorials) < TutorialsPerPagination {
		tutSlice = tutorials
	} else {
		tutSlice = tutorials[:len(tutorials)-1]
		lastTut = tutorials[len(tutorials)-1]
	}

	pageData.Tutorials = &html.TutorialsListComponent{
		Tutorials:    tutSlice,
		LastTutorial: lastTut,
		QueryURL:     fmt.Sprintf("/tutorials/page/%d", 2),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "tutorials", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialsPaginationGet(w http.ResponseWriter, r *http.Request) {
	tutorialsComponent := &html.TutorialsListComponent{}

	pageNumber, err := strconv.Atoi(chi.URLParam(r, "page-number"))
	if err != nil {
		h.ErrorLog.Printf("Failed to convert page-number to int: %s\n", err)
		tutorialsComponent.ErrorMessage = "Unexpected server error. Please try again."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "tutorials", tutorialsComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	tutorials, err := h.Database.GetAllTutorialsPaginated(pageNumber, TutorialsPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all tutorials (page %d) from the database: %s\n", pageNumber, err)
		tutorialsComponent.ErrorMessage = "Failed to fetch next tutorials."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "tutorials", tutorialsComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var tutSlice []*models.TutorialModel
	var lastTut *models.TutorialModel

	if len(tutorials) < TutorialsPerPagination {
		tutSlice = tutorials
	} else {
		tutSlice = tutorials[:len(tutorials)-1]
		lastTut = tutorials[len(tutorials)-1]
	}

	tutorialsComponent.Tutorials = tutSlice
	tutorialsComponent.LastTutorial = lastTut
	tutorialsComponent.QueryURL = fmt.Sprintf("/tutorials/page/%d", pageNumber+1)

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "tutorials", tutorialsComponent); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialsSearchGet(w http.ResponseWriter, r *http.Request) {
	var page int

	queryPage := r.URL.Query().Get("page")
	query := r.URL.Query().Get("query")
	tutorialsComponent := &html.TutorialsListComponent{}

	if queryPage == "" {
		page = 1
	} else {
		pageNum, err := strconv.Atoi(queryPage)
		if err != nil {
			h.WarningLog.Printf("Failed to convert page to int: %s\n", err)
			page = 1
		} else {
			page = pageNum
		}
	}

	tutorials, err := h.Database.SearchTutorialsPaginated(query, page, TutorialsPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to search for tutorials (page %d) from the database: %s\n", page, err)
		tutorialsComponent.ErrorMessage = "Failed to get tutorials. Please try again."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "tutorials", tutorialsComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var tutSlice []*models.TutorialModel
	var lastTut *models.TutorialModel

	if len(tutorials) < TutorialsPerPagination {
		tutSlice = tutorials
	} else {
		tutSlice = tutorials[:len(tutorials)-1]
		lastTut = tutorials[len(tutorials)-1]
	}

	tutorialsComponent.Tutorials = tutSlice
	tutorialsComponent.LastTutorial = lastTut
	tutorialsComponent.QueryURL = fmt.Sprintf("/tutorials/search?page=%d&query=%s", page+1, url.QueryEscape(query))

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "tutorials", tutorialsComponent); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.TutorialsTutorialPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "tutorials-tutorial", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}
