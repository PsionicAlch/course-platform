package comments

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

const CommentsPerPagination = 25

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
		Renderers: *pages.CreateRenderers(pageRenderer, htmxRenderer),
		Database:  db,
		Auth:      auth,
	}
}

func (h *Handlers) CommentsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminCommentsPage{
		BasePage: html.NewBasePage(user),
	}

	commentsList, urlQuery, err := h.CreateCommentsList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create comments list: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Comments = commentsList

	urlQuery.Set("page", "1")
	pageData.URLQuery = urlQuery.Encode()

	tutorials, err := h.Database.GetAllTutorials()
	if err != nil {
		h.ErrorLog.Printf("Failed to get all tutorials: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Tutorials = tutorials

	users, err := h.Database.GetAllUsers()
	if err != nil {
		h.ErrorLog.Printf("Failed to get all users: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Users = users

	numComments, err := h.Database.CountComments()
	if err != nil {
		h.ErrorLog.Printf("Failed to count the omments: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.NumComments = numComments

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-comments", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CommentsPaginationGet(w http.ResponseWriter, r *http.Request) {
	commentsList, _, err := h.CreateCommentsList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create comments list: %s\n", err)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "admin-comments", html.AdminCommentsListComponent{ErrorMessage: "Failed to get comments. Please try again."}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "admin-comments", commentsList); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CommentDelete(w http.ResponseWriter, r *http.Request) {
	urlQuery := CreateUrlQuery(r)

	commentId := chi.URLParam(r, "comment-id")

	if err := h.Database.DeleteComment(commentId); err != nil {
		h.ErrorLog.Printf("Failed to delete comment (\"%s\"): %s\n", commentId, err)

		resp := fmt.Sprintf(`<button class="btn btn-red shadow-sm" hx-delete="/admin/comments/%s">Delete Comment</button>`, commentId)
		resp += `
		<script>
		notyf.open({
			type: 'flash-error',
			message: 'Failed to delete comment. Please try again.'
		});
		</script>
		`

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/admin/comments?%s", urlQuery.Encode()))
}

// Possible URL queries:
// - page
// - query
// - tutorial
// - user
func (h *Handlers) CreateCommentsList(r *http.Request) (*html.AdminCommentsListComponent, url.Values, error) {
	urlQuery := CreateUrlQuery(r)

	query := urlQuery.Get("query")
	tutorialId := urlQuery.Get("tutorial")
	userId := urlQuery.Get("user")
	page := 1

	if p, err := strconv.Atoi(urlQuery.Get("page")); err == nil {
		page = p - 1
	}

	comments, err := h.Database.AdminGetComments(query, tutorialId, userId, uint(page), CommentsPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all comments (page %d) from the database: %s\n", page, err)
		return nil, urlQuery, err
	}

	var commentsSlice []*models.CommentModel
	var lastComment *models.CommentModel

	if len(comments) < CommentsPerPagination {
		commentsSlice = comments
	} else {
		commentsSlice = comments[:len(comments)-1]
		lastComment = comments[len(comments)-1]
	}

	users := make(map[string]*models.UserModel, len(comments))
	tutorials := make(map[string]*models.TutorialModel, len(comments))

	for _, comment := range comments {
		user, err := h.Database.GetUserByID(comment.UserID, database.All)
		if err != nil {
			h.ErrorLog.Printf("Failed to get user (\"%s\") by ID: %s\n", comment.UserID, err)
			return nil, urlQuery, err
		}

		users[comment.ID] = user

		tutorial, err := h.Database.GetTutorialByID(comment.TutorialID)
		if err != nil {
			h.ErrorLog.Printf("Failed to get tutorial (\"%s\") by ID: %s\n", comment.TutorialID, err)
			return nil, urlQuery, err
		}

		tutorials[comment.ID] = tutorial
	}

	commentsList := &html.AdminCommentsListComponent{
		Comments:    commentsSlice,
		LastComment: lastComment,
		Tutorials:   tutorials,
		Users:       users,
		BaseURL:     "/admin/comments/htmx",
		URLQuery:    urlQuery.Encode(),
	}

	return commentsList, urlQuery, nil
}

func CreateUrlQuery(r *http.Request) url.Values {
	urlQuery := make(url.Values)

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		urlQuery.Add("page", fmt.Sprintf("%d", p+1))
	} else {
		urlQuery.Add("page", "2")
	}

	if q := r.URL.Query().Get("query"); q != "" {
		urlQuery.Add("query", q)
	}

	if t := r.URL.Query().Get("tutorial"); t != "" {
		urlQuery.Add("tutorial", t)
	}

	if u := r.URL.Query().Get("user"); u != "" {
		urlQuery.Add("user", u)
	}

	return urlQuery
}
