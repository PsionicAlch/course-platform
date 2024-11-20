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
	Auth      *authentication.Authentication
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, db database.Database, sessions *session.Session, auth *authentication.Authentication) *Handlers {
	loggers := utils.CreateLoggers("TUTORIALS HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: pages.CreateRenderers(pageRenderer, htmxRenderer),
		Database:  db,
		Session:   sessions,
		Auth:      auth,
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
		User:     user,
	}

	tutorialSlug := chi.URLParam(r, "slug")

	tutorial, err := h.Database.GetTutorialBySlug(tutorialSlug)
	if err != nil {
		h.ErrorLog.Printf("Failed to get tutorial (\"%s\") in the database: %s\n", tutorialSlug, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if tutorial == nil {
		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user)}, http.StatusNotFound); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if user != nil {
		userLikedTutorial, err := h.Database.UserLikedTutorial(user.ID, tutorialSlug)
		if err != nil {
			h.ErrorLog.Printf("Failed to find out if user liked tutorial (\"%s\") from the database: %s\n", tutorialSlug, err)

			if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		userBookmarkedTutorial, err := h.Database.UserBookmarkedTutorial(user.ID, tutorialSlug)
		if err != nil {
			h.ErrorLog.Printf("Failed to find out if user bookmarked tutorial (\"%s\") from the database: %s\n", tutorialSlug, err)

			if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		pageData.TutorialLiked = userLikedTutorial
		pageData.TutorialBookmarked = userBookmarkedTutorial
	}

	pageData.Tutorial = tutorial
	pageData.Author = &models.AuthorModel{
		Name:    "Jean-Jacques",
		Surname: "Strydom",
		Slug:    "jean-jacques-strydom",
	}
	pageData.Course = nil

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "tutorials-tutorial", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) LikeTutorialPost(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	tutorialSlug := chi.URLParam(r, "slug")

	if user == nil {
		if err := h.Renderers.Htmx.RenderHTML(w, nil, "heart-empty", nil); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	likedTutorial, err := h.Database.UserLikedTutorial(user.ID, tutorialSlug)
	if err != nil {
		h.ErrorLog.Printf("Failed to see if user liked tutorial %s in the database: %s\n", tutorialSlug, err)
		if err := h.Renderers.Htmx.RenderHTML(w, nil, "heart-empty", nil, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if likedTutorial {
		if err := h.Database.UserDislikeTutorial(user.ID, tutorialSlug); err != nil {
			h.ErrorLog.Printf("Failed to dislike tutorial %s in the database: %s\n", tutorialSlug, err)
			if err := h.Renderers.Htmx.RenderHTML(w, nil, "heart-empty", nil, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "heart-empty", nil); err != nil {
			h.ErrorLog.Println(err)
		}
	} else {
		if err := h.Database.UserLikeTutorial(user.ID, tutorialSlug); err != nil {
			h.ErrorLog.Printf("Failed to like tutorial %s in the database: %s\n", tutorialSlug, err)
			if err := h.Renderers.Htmx.RenderHTML(w, nil, "heart-filled", nil, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "heart-filled", nil); err != nil {
			h.ErrorLog.Println(err)
		}
	}
}

func (h *Handlers) BookmarkTutorialPost(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	tutorialSlug := chi.URLParam(r, "slug")

	if user == nil {
		if err := h.Renderers.Htmx.RenderHTML(w, nil, "bookmark-empty", nil); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	bookmarkedTutorial, err := h.Database.UserBookmarkedTutorial(user.ID, tutorialSlug)
	if err != nil {
		h.ErrorLog.Printf("Failed to see if user bookmarked tutorial %s in the database: %s\n", tutorialSlug, err)
		if err := h.Renderers.Htmx.RenderHTML(w, nil, "bookmark-empty", nil, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if bookmarkedTutorial {
		if err := h.Database.UserUnbookmarkTutorial(user.ID, tutorialSlug); err != nil {
			h.ErrorLog.Printf("Failed to unbookmark tutorial %s in the database: %s\n", tutorialSlug, err)
			if err := h.Renderers.Htmx.RenderHTML(w, nil, "bookmark-empty", nil, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "bookmark-empty", nil); err != nil {
			h.ErrorLog.Println(err)
		}
	} else {
		if err := h.Database.UserBookmarkTutorial(user.ID, tutorialSlug); err != nil {
			h.ErrorLog.Printf("Failed to bookmark tutorial %s in the database: %s\n", tutorialSlug, err)
			if err := h.Renderers.Htmx.RenderHTML(w, nil, "bookmark-filled", nil, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "bookmark-filled", nil); err != nil {
			h.ErrorLog.Println(err)
		}
	}
}
