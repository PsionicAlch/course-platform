package courses

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

const CoursesPerPagination = 25

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
		Renderers: *pages.CreateRenderers(pageRenderer, htmxRenderer, nil),
		Database:  db,
		Auth:      auth,
	}
}

func (h *Handlers) CoursesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminCoursesPage{
		BasePage: html.NewBasePage(user),
	}

	coursesList, urlQuery, err := h.CreateCoursesList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create courses list: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Courses = coursesList

	urlQuery.Set("page", "1")
	pageData.URLQuery = urlQuery.Encode()

	pageData.PublishStatus = PublishStatuses

	numCourses, err := h.Database.CountCourses()
	if err != nil {
		h.ErrorLog.Printf("Failed to count the number of courses in the database: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.NumCourses = numCourses

	authors, err := h.Database.GetUsers("", database.Author, "", "")
	if err != nil {
		h.ErrorLog.Printf("Failed to get all the authors from the database: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Authors = authors

	keywords, err := h.Database.GetKeywords()
	if err != nil {
		h.ErrorLog.Printf("Failed to get all the keywords from the database: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Keywords = keywords

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-courses", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CoursesPaginationGet(w http.ResponseWriter, r *http.Request) {
	coursesList, _, err := h.CreateCoursesList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create courses list component: %s\n", err)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "admin-courses", html.AdminCoursesListComponent{ErrorMessage: "Failed to load courses. Please try again."}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "admin-courses", coursesList); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PublishedEditGet(w http.ResponseWriter, r *http.Request) {
	courseId := chi.URLParam(r, "course-id")

	course, err := h.Database.GetCourseByID(courseId)
	if err != nil {
		h.ErrorLog.Printf("Failed to get course by ID \"%s\": %s\n", courseId, err)

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

	if course == nil {
		h.ErrorLog.Printf("Failed to get course by ID \"%s\": Nil was returned\n", courseId)

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
	if course.Published {
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
		URL:      fmt.Sprintf("/admin/courses/htmx/change-published/%s", course.ID),
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", selectComponent); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PublishedEditPost(w http.ResponseWriter, r *http.Request) {
	courseId := chi.URLParam(r, "course-id")

	r.ParseForm()
	publishStatus := r.Form.Get("publish-status")

	course, err := h.Database.GetCourseByID(courseId)
	if err != nil {
		h.ErrorLog.Printf("Failed to get course by ID \"%s\": %s\n", courseId, err)

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

	if course == nil {
		h.ErrorLog.Printf("Failed to get course by ID \"%s\": Nil was returned\n", courseId)

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
			URL:          fmt.Sprintf("/admin/courses/change-published/%s", course.ID),
			ErrorMessage: "Invalid publish status selected.",
		}

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", selectComponent, http.StatusBadRequest); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if publishStatus == "Published" {
		if err := h.Database.PublishCourse(course.ID); err != nil {
			h.ErrorLog.Printf("Failed to update tutorial's (\"%s\") publish status: %s\n", course.Title, err)

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
		if err := h.Database.UnpublishCourse(course.ID); err != nil {
			h.ErrorLog.Printf("Failed to update course's (\"%s\") publish status: %s\n", course.Title, err)

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
	courseId := chi.URLParam(r, "course-id")

	course, err := h.Database.GetCourseByID(courseId)
	if err != nil {
		h.ErrorLog.Printf("Failed to get course by ID \"%s\": %s\n", courseId, err)

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

	authors, err := h.Database.GetUsers("", database.Author, "", "")
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

	if course.AuthorID.Valid {
		selected = course.AuthorID.String
	}

	selectComponent := html.SelectComponent{
		Name:     "author",
		Options:  selectOptions,
		Selected: selected,
		URL:      fmt.Sprintf("/admin/courses/htmx/change-author/%s", course.ID),
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", selectComponent); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) AuthorEditPost(w http.ResponseWriter, r *http.Request) {
	courseId := chi.URLParam(r, "course-id")

	r.ParseForm()
	authorId := r.Form.Get("author")

	course, err := h.Database.GetCourseByID(courseId)
	if err != nil {
		h.ErrorLog.Printf("Failed to get course by ID \"%s\": %s\n", courseId, err)

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

	if err := h.Database.UpdateCourseAuthor(course.ID, authorId); err != nil {
		h.ErrorLog.Printf("Failed to update course's \"%s\" author \"%s\": %s\n", course.ID, authorId, err)

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
			h.ErrorLog.Printf("Failed to update course's \"%s\" author \"%s\": %s\n", course.ID, authorId, err)

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
// -keyword
// -bought_by
func (h *Handlers) CreateCoursesList(r *http.Request) (*html.AdminCoursesListComponent, url.Values, error) {
	var published *bool
	var page int
	var query string
	var author *string
	var boughtBy string
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

	if bought := r.URL.Query().Get("bought_by"); bought != "" {
		boughtBy = bought

		urlQuery.Add("bought_by", bought)
	}

	if key := r.URL.Query().Get("keyword"); key != "" {
		keyword = key

		urlQuery.Add("keyword", key)
	}

	courses, err := h.Database.AdminGetCourses(query, published, author, boughtBy, keyword, uint(page), CoursesPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get courses (page %d) from the database: %s\n", page, err)
		return nil, urlQuery, err
	}

	var coursesSlice []*models.CourseModel
	var lastCourse *models.CourseModel

	if len(courses) < CoursesPerPagination {
		coursesSlice = courses
	} else {
		coursesSlice = courses[:len(courses)-1]
		lastCourse = courses[len(courses)-1]
	}

	authors := make(map[string]*models.UserModel, len(courses))
	keywords := make(map[string][]string, len(courses))
	purchases := make(map[string]uint, len(courses))

	for _, course := range courses {
		if course.AuthorID.Valid {
			author, err := h.Database.GetUserByID(course.AuthorID.String, database.Author)
			if err != nil {
				h.ErrorLog.Printf("Failed to get author \"%s\" from the database: %s\n", course.AuthorID.String, err)
				return nil, urlQuery, err
			}

			authors[course.ID] = author
		}

		keys, err := h.Database.GetAllKeywordsForCourse(course.ID)
		if err != nil {
			h.ErrorLog.Printf("Failed to get all keywords for course \"%s\": %s\n", course.Title, err)
			return nil, urlQuery, err
		}

		keywords[course.ID] = keys

		purchase, err := h.Database.CountUsersWhoBoughtCourse(course.ID)
		if err != nil {
			h.ErrorLog.Printf("Failed to count the number of times the course \"%s\" was bought: %s\n", course.Title, err)
			return nil, urlQuery, err
		}

		purchases[course.ID] = purchase
	}

	coursesList := &html.AdminCoursesListComponent{
		Courses:    coursesSlice,
		LastCourse: lastCourse,
		Authors:    authors,
		Keywords:   keywords,
		Purchases:  purchases,
		BaseURL:    "/admin/courses/htmx",
		URLQuery:   urlQuery.Encode(),
	}

	return coursesList, urlQuery, nil
}
