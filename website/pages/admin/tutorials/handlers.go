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

var PublishStatuses = []string{"Published", "Unpublished"}

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

func (h *Handlers) TutorialsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminTutorialsPage{
		BasePage: html.NewBasePage(user),
	}

	tutorialList, urlQuery, err := h.CreateTutorialsList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to construct tutorial list component: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Tutorials = tutorialList

	urlQuery.Set("page", "1")

	pageData.URLQuery = urlQuery.Encode()

	pageData.PublishStatus = PublishStatuses

	numTutorials, err := h.Database.CountTutorials()
	if err != nil {
		h.ErrorLog.Printf("Failed to count the number of tutorials in the database: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.NumTutorials = numTutorials

	authors, err := h.Database.GetUsers("", database.Author)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all the authors from the database: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Authors = authors

	keywords, err := h.Database.GetKeywords()
	if err != nil {
		h.ErrorLog.Printf("Failed to get all the keywords from the database: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Keywords = keywords

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-tutorials", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialsPaginationGet(w http.ResponseWriter, r *http.Request) {
	tutorialsList, _, err := h.CreateTutorialsList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create tutorials list component: %s\n", err)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "admin-tutorials", html.AdminTutorialsListComponent{ErrorMessage: "Failed to load tutorials. Please try again."}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "admin-tutorials", tutorialsList); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PublishedEditGet(w http.ResponseWriter, r *http.Request) {
	tutorialId := chi.URLParam(r, "tutorial-id")

	tutorial, err := h.Database.GetTutorialByID(tutorialId)
	if err != nil {
		h.ErrorLog.Printf("Failed to get tutorial by ID \"%s\": %s\n", tutorialId, err)

		resp := "Unpublished"
		resp += `
		<script>
			notyf.open({
				type: 'flash-error',
				message: "Unexpected server error"
			});
		</script>
		`

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if tutorial == nil {
		h.ErrorLog.Printf("Failed to get tutorial by ID \"%s\": Nil was returned\n", tutorialId)

		resp := "Unpublished"
		resp += `
		<script>
			notyf.open({
				type: 'flash-error',
				message: "Unexpected server error"
			});
		</script>
		`

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var publishStatus string
	if tutorial.Published {
		publishStatus = "Published"
	} else {
		publishStatus = "Unpublished"
	}

	publishStatuses := make(map[string]string, len(PublishStatuses))
	for _, status := range PublishStatuses {
		publishStatuses[status] = status
	}

	selectComponent := html.SelectComponent{
		Name:     "publish-status",
		Options:  publishStatuses,
		Selected: publishStatus,
		URL:      fmt.Sprintf("/admin/tutorials/htmx/change-published/%s", tutorialId),
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", selectComponent); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PublishedEditPost(w http.ResponseWriter, r *http.Request) {
	tutorialId := chi.URLParam(r, "tutorial-id")

	r.ParseForm()
	publishStatus := r.Form.Get("publish-status")

	tutorial, err := h.Database.GetTutorialByID(tutorialId)
	if err != nil {
		h.ErrorLog.Printf("Failed to get tutorial by ID \"%s\": %s\n", tutorialId, err)

		resp := "Unpublished"
		resp += `
		<script>
			notyf.open({
				type: 'flash-error',
				message: "Unexpected server error"
			});
		</script>
		`

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if tutorial == nil {
		h.ErrorLog.Printf("Failed to get tutorial by ID \"%s\": Nil was returned\n", tutorialId)

		resp := "Unpublished"
		resp += `
		<script>
			notyf.open({
				type: 'flash-error',
				message: "Unexpected server error"
			});
		</script>
		`

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if !utils.InSlice(publishStatus, PublishStatuses) {
		publishStatuses := make(map[string]string, len(PublishStatuses))
		for _, status := range PublishStatuses {
			publishStatuses[status] = status
		}

		selectComponent := html.SelectComponent{
			Name:         "publish-status",
			Options:      publishStatuses,
			URL:          fmt.Sprintf("/admin/tutorials/change-published/%s", tutorialId),
			ErrorMessage: "Invalid publish status selected.",
		}

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", selectComponent, http.StatusBadRequest); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if publishStatus == "Published" {
		if err := h.Database.PublishTutorial(tutorial.ID); err != nil {
			h.ErrorLog.Printf("Failed to update tutorial's (\"%s\") publish status: %s\n", tutorial.Title, err)

			resp := "Unpublished"
			resp += `
            <script>
                notyf.open({
                    type: 'flash-error',
                    message: "Unexpected server error"
                });
            </script>
            `

			if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", "Published"); err != nil {
			h.ErrorLog.Println(err)
		}
	} else {
		if err := h.Database.PublishTutorial(tutorial.ID); err != nil {
			h.ErrorLog.Printf("Failed to update tutorial's (\"%s\") publish status: %s\n", tutorial.Title, err)

			resp := "Published"
			resp += `
            <script>
                notyf.open({
                    type: 'flash-error',
                    message: "Unexpected server error"
                });
            </script>
            `

			if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", "Unpublished"); err != nil {
			h.ErrorLog.Println(err)
		}
	}
}

func (h *Handlers) AuthorEditGet(w http.ResponseWriter, r *http.Request) {
	tutorialId := chi.URLParam(r, "tutorial-id")

	tutorial, err := h.Database.GetTutorialByID(tutorialId)
	if err != nil {
		h.ErrorLog.Printf("Failed to get tutorial by ID \"%s\": %s\n", tutorialId, err)

		resp := "No Author"
		resp += `
		<script>
			notyf.open({
				type: 'flash-error',
				message: "Unexpected server error"
			});
		</script>
		`

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	authors, err := h.Database.GetUsers("", database.Author)
	if err != nil {
		h.ErrorLog.Printf("Failed to get authors: %s\n", err)

		resp := "No Author"
		resp += `
		<script>
			notyf.open({
				type: 'flash-error',
				message: "Unexpected server error"
			});
		</script>
		`

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	selectOptions := map[string]string{
		"": "No Author",
	}

	for _, author := range authors {
		selectOptions[author.ID] = fmt.Sprintf("%s %s", author.Name, author.Surname)
	}

	var selected string

	if tutorial.AuthorID.Valid {
		selected = tutorial.AuthorID.String
	}

	selectComponent := html.SelectComponent{
		Name:     "author",
		Options:  selectOptions,
		Selected: selected,
		URL:      fmt.Sprintf("/admin/tutorials/htmx/change-author/%s", tutorial.ID),
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", selectComponent); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) AuthorEditPost(w http.ResponseWriter, r *http.Request) {
	tutorialId := chi.URLParam(r, "tutorial-id")

	r.ParseForm()

	authorId := r.Form.Get("author")

	tutorial, err := h.Database.GetTutorialByID(tutorialId)
	if err != nil {
		h.ErrorLog.Printf("Failed to get tutorial by ID \"%s\": %s\n", tutorialId, err)

		resp := "No Author"
		resp += `
		<script>
			notyf.open({
				type: 'flash-error',
				message: "Unexpected server error"
			});
		</script>
		`

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Database.UpdateAuthor(tutorial.ID, authorId); err != nil {
		h.ErrorLog.Printf("Failed to update tutorial \"%s\" author \"%s\": %s\n", tutorial.ID, authorId, err)

		resp := "No Author"
		resp += `
		<script>
			notyf.open({
				type: 'flash-error',
				message: "Unexpected server error"
			});
		</script>
		`

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var resp string

	if authorId == "" {
		resp = "No Author"
	} else {
		author, err := h.Database.GetUserByID(authorId, database.Author)
		if err != nil {
			h.ErrorLog.Printf("Failed to update tutorial \"%s\" author \"%s\": %s\n", tutorial.ID, authorId, err)

			resp := "No Author"
			resp += `
            <script>
                notyf.open({
                    type: 'flash-error',
                    message: "Unexpected server error"
                });
            </script>
            `

			if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		resp = fmt.Sprintf("%s %s", author.Name, author.Surname)
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp); err != nil {
		h.ErrorLog.Println(err)
	}
}

// Possible URL queries:
// -page
// -query
// -status
// -author
// -liked_by
// -bookmarked_by
// -keyword
func (h *Handlers) CreateTutorialsList(r *http.Request) (*html.AdminTutorialsListComponent, url.Values, error) {
	var published *bool
	var page int
	var query string
	var author *string
	var likedBy string
	var bookmarkedBy string
	var keyword string

	urlQuery := make(url.Values)

	if !utils.InSlice(r.URL.Query().Get("status"), PublishStatuses) {
		published = nil
	} else {
		if r.URL.Query().Get("status") == "Published" {
			tmp := true
			published = &tmp
		} else {
			tmp := false
			published = &tmp
		}

		urlQuery.Add("status", r.URL.Query().Get("status"))
	}

	if pageNum, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && pageNum > 0 {
		page = pageNum
	} else {
		page = 1
	}

	urlQuery.Add("page", strconv.Itoa(page+1))

	if r.URL.Query().Get("query") != "" {
		query = r.URL.Query().Get("query")
		urlQuery.Add("query", query)
	}

	if authorStr := r.URL.Query().Get("author"); authorStr != "" {
		if authorStr == "nil" {
			author = nil
		} else {
			author = &authorStr
		}

		urlQuery.Add("author", authorStr)
	} else {
		temp := ""
		author = &temp
	}

	if liked := r.URL.Query().Get("liked_by"); liked != "" {
		likedBy = liked

		urlQuery.Add("liked_by", liked)
	}

	if bookmarked := r.URL.Query().Get("bookmarked_by"); bookmarked != "" {
		bookmarkedBy = bookmarked

		urlQuery.Add("bookmarked_by", bookmarked)
	}

	if key := r.URL.Query().Get("keyword"); key != "" {
		keyword = key

		urlQuery.Add("keyword", key)
	}

	tutorials, err := h.Database.GetTutorials(query, published, author, likedBy, bookmarkedBy, keyword, uint(page), TutorialsPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get the tutorials from the database")
		return nil, urlQuery, err
	}

	var tutorialsSlice []*models.TutorialModel
	var lastTutorial *models.TutorialModel

	if len(tutorials) < TutorialsPerPagination {
		tutorialsSlice = tutorials
	} else {
		tutorialsSlice = tutorials[:len(tutorials)-1]
		lastTutorial = tutorials[len(tutorials)-1]
	}

	authors := make(map[string]*models.UserModel, len(tutorials))
	keywords := make(map[string][]string, len(tutorials))
	comments := make(map[string]uint, len(tutorials))
	likes := make(map[string]uint, len(tutorials))
	bookmarks := make(map[string]uint, len(tutorials))

	for _, tutorial := range tutorials {
		if tutorial.AuthorID.Valid {
			author, err := h.Database.GetUserByID(tutorial.AuthorID.String, database.Author)
			if err != nil {
				h.ErrorLog.Printf("Failed to get author \"%s\" from the database: %s\n", tutorial.AuthorID.String, err)
				return nil, urlQuery, err
			}

			authors[tutorial.ID] = author
		}

		keys, err := h.Database.GetAllKeywordsForTutorial(tutorial.ID)
		if err != nil {
			h.ErrorLog.Printf("Failed to get all keywords for tutorial \"%s\": %s\n", tutorial.Title, err)
			return nil, urlQuery, err
		}

		keywords[tutorial.ID] = keys

		commentCount, err := h.Database.CountCommentsForTutorial(tutorial.ID)
		if err != nil {
			h.ErrorLog.Printf("Failed to count all comments related to tutorial \"%s\": %s\n", tutorial.Title, err)
			return nil, urlQuery, err
		}

		comments[tutorial.ID] = commentCount

		likesCount, err := h.Database.CountTutorialLikes(tutorial.ID)
		if err != nil {
			h.ErrorLog.Printf("Failed to count the number of likes the tutorial \"%s\" has: %s\n", tutorial.Title, err)
			return nil, urlQuery, err
		}

		likes[tutorial.ID] = likesCount

		bookmarksCount, err := h.Database.CountTutorialBookmarks(tutorial.ID)
		if err != nil {
			h.ErrorLog.Printf("Failed to count the number of bookmarks the tutorial \"%s\" has: %s\n", tutorial.Title, err)
			return nil, urlQuery, err
		}

		bookmarks[tutorial.ID] = bookmarksCount
	}

	usersList := &html.AdminTutorialsListComponent{
		Tutorials:    tutorialsSlice,
		LastTutorial: lastTutorial,
		Authors:      authors,
		Keywords:     keywords,
		Comments:     comments,
		Likes:        likes,
		Bookmarks:    bookmarks,
		BaseURL:      "/admin/tutorials/htmx",
		URLQuery:     urlQuery.Encode(),
	}

	return usersList, urlQuery, nil
}
