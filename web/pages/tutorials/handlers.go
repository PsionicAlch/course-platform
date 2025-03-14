package tutorials

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PsionicAlch/course-platform/internal/authentication"
	"github.com/PsionicAlch/course-platform/internal/database"
	"github.com/PsionicAlch/course-platform/internal/database/models"
	"github.com/PsionicAlch/course-platform/internal/utils"
	"github.com/PsionicAlch/course-platform/web/html"
	"github.com/PsionicAlch/course-platform/web/pages"
	goaway "github.com/TwiN/go-away"
	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
)

const TutorialsPerPagination = 25
const CommentsPerPagination = 25

type Handlers struct {
	utils.Loggers
	*pages.HandlerContext
}

func SetupHandlers(handlerContext *pages.HandlerContext) *Handlers {
	loggers := utils.CreateLoggers("TUTORIALS HANDLERS")

	return &Handlers{
		Loggers:        loggers,
		HandlerContext: handlerContext,
	}
}

func (h *Handlers) TutorialsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.TutorialsPage{
		BasePage: html.NewBasePage(user, nosurf.Token(r)),
	}

	tutorialsList, err := h.CreateTutorialsList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create tutorials list: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}
		return
	}

	pageData.Tutorials = tutorialsList

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "tutorials", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialsPaginationGet(w http.ResponseWriter, r *http.Request) {
	tutorialsList, err := h.CreateTutorialsList(r)
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

func (h *Handlers) TutorialGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.TutorialsTutorialPage{
		BasePage: html.NewBasePage(user, nosurf.Token(r)),
		User:     user,
	}

	tutorialSlug := chi.URLParam(r, "slug")

	tutorial, err := h.Database.GetTutorialBySlug(tutorialSlug)
	if err != nil {
		h.ErrorLog.Printf("Failed to get tutorial (\"%s\") in the database: %s\n", tutorialSlug, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if tutorial == nil {
		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusNotFound); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Tutorial = tutorial

	keywords, err := h.Database.GetAllKeywordsForTutorial(tutorial.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to get keywords for tutorial (\"%s\") in the database: %s\n", tutorial.Title, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Keywords = keywords

	if user != nil {
		userLikedTutorial, err := h.Database.UserLikedTutorial(user.ID, tutorialSlug)
		if err != nil {
			h.ErrorLog.Printf("Failed to find out if user liked tutorial (\"%s\") from the database: %s\n", tutorialSlug, err)

			if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		userBookmarkedTutorial, err := h.Database.UserBookmarkedTutorial(user.ID, tutorialSlug)
		if err != nil {
			h.ErrorLog.Printf("Failed to find out if user bookmarked tutorial (\"%s\") from the database: %s\n", tutorialSlug, err)

			if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		pageData.TutorialLiked = userLikedTutorial
		pageData.TutorialBookmarked = userBookmarkedTutorial
	}

	var authorId string
	if tutorial.AuthorID.Valid {
		authorId = tutorial.AuthorID.String
	}

	author, err := h.Database.GetUserByID(authorId, database.Author)
	if err != nil {
		h.ErrorLog.Printf("Failed to get author by ID (\"%s\") from the database: %s\n", authorId, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Author = author

	var course *models.CourseModel

	courses, err := h.Database.GetCourses("", "", 1, 1)
	if err != nil {
		h.ErrorLog.Printf("Failed to get author by ID (\"%s\") from the database: %s\n", authorId, err)
	}

	if len(courses) >= 1 {
		course = courses[0]
	}

	pageData.Course = course

	comments, err := h.Database.GetAllCommentsPaginated(tutorial.ID, 1, CommentsPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get comments for tutorial (\"%s\"): %s\n", tutorial.Title, err)
		h.Session.SetErrorMessage(r.Context(), "Failed to get comments for tutorial.")
	}

	users := make(map[string]*models.UserModel)

	for _, comment := range comments {
		user, err := h.Database.GetUserByID(comment.UserID, database.All)
		if err != nil {
			h.ErrorLog.Printf("Failed to get user (\"%s\") for comment (\"%s\"): %s\n", comment.ID, comment.UserID, err)
			h.Session.SetErrorMessage(r.Context(), "Failed to get comments for tutorial.")
		}

		users[comment.UserID] = user
	}

	var commentsSlice []*models.CommentModel
	var lastComment *models.CommentModel

	if len(comments) < CommentsPerPagination {
		commentsSlice = comments
	} else {
		commentsSlice = comments[:len(comments)-1]
		lastComment = comments[len(comments)-1]
	}

	pageData.Comments = &html.CommentsListComponent{
		Comments:    commentsSlice,
		LastComment: lastComment,
		Users:       users,
		QueryURL:    fmt.Sprintf("/tutorials/%s/comments?page=2", tutorial.Slug),
	}

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

func (h *Handlers) CommentsGet(w http.ResponseWriter, r *http.Request) {
	commentsList := &html.CommentsListComponent{}

	tutorialSlug := chi.URLParam(r, "slug")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		h.ErrorLog.Printf("Failed to convert page URL query to int: %s\n", err)
		commentsList.ErrorMessage = "Failed to load comments."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "comments", commentsList, http.StatusBadRequest); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	comments, err := h.Database.GetAllCommentsBySlugPaginated(tutorialSlug, page, CommentsPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get comments for tutorial (\"%s\"): %s\n", tutorialSlug, err)
		commentsList.ErrorMessage = "Failed to load comments."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "comments", commentsList, http.StatusBadRequest); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	users := make(map[string]*models.UserModel)

	for _, comment := range comments {
		user, err := h.Database.GetUserByID(comment.UserID, database.All)
		if err != nil {
			h.ErrorLog.Printf("Failed to get user (\"%s\") for comment (\"%s\"): %s\n", comment.ID, comment.UserID, err)
			commentsList.ErrorMessage = "Failed to load comments."

			if err := h.Renderers.Htmx.RenderHTML(w, nil, "comments", commentsList, http.StatusBadRequest); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		users[comment.UserID] = user
	}

	var commentsSlice []*models.CommentModel
	var lastComment *models.CommentModel

	if len(comments) < CommentsPerPagination {
		commentsSlice = comments
	} else {
		commentsSlice = comments[:len(comments)-1]
		lastComment = comments[len(comments)-1]
	}

	commentsList.Comments = commentsSlice
	commentsList.LastComment = lastComment
	commentsList.Users = users
	commentsList.QueryURL = fmt.Sprintf("/tutorials/%s/comments?page=%d", tutorialSlug, page+1)

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "comments", commentsList); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CommentsPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	tutorialSlug := chi.URLParam(r, "slug")
	comment := r.Form.Get("comment")
	user := authentication.GetUserFromRequest(r)

	if user != nil && comment != "" {
		if len(comment) > 500 {
			comment = comment[:500]
		}

		comment = goaway.Censor(comment)

		commentModel, err := h.Database.AddCommentBySlug(comment, user.ID, tutorialSlug)
		if err != nil {
			h.ErrorLog.Printf("Failed to add user's (\"%s\") comment on \"%s\" to the database: %s\n", user.ID, tutorialSlug, err)
			return
		}

		commentsList := &html.CommentsListComponent{}
		commentsList.Comments = []*models.CommentModel{commentModel}
		commentsList.Users = map[string]*models.UserModel{commentModel.UserID: user}

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "comments", commentsList); err != nil {
			h.ErrorLog.Println(err)
		}
	}
}

// Possible URL queries:
// -page
// -query
func (h *Handlers) CreateTutorialsList(r *http.Request) (*html.TutorialsListComponent, error) {
	query := r.URL.Query().Get("query")
	page := 1

	urlQuery := make(url.Values)

	if pageNum, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		page = pageNum
	}

	urlQuery.Add("query", query)
	urlQuery.Add("page", fmt.Sprintf("%d", page+1))

	tutorials, err := h.Database.GetTutorials(query, "", page, TutorialsPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all tutorials (page %d): %s\n", page, err)
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

	tutorialList := &html.TutorialsListComponent{
		Tutorials:    tutorialsSlice,
		LastTutorial: lastTutorial,
		QueryURL:     fmt.Sprintf("/tutorials/htmx?%s", urlQuery.Encode()),
	}

	return tutorialList, nil
}
