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
	"github.com/justinas/nosurf"
)

const CoursesPerPagination = 25

type Handlers struct {
	utils.Loggers
	Renderers pages.Renderers
	Database  database.Database
}

func SetupHandlers(pageRenderer, htmxRenderer render.Renderer, db database.Database) *Handlers {
	loggers := utils.CreateLoggers("AUTHOR COURSES HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: *pages.CreateRenderers(pageRenderer, htmxRenderer, nil),
		Database:  db,
	}
}

func (h *Handlers) CoursesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AuthorsCoursesPage{
		BasePage: html.NewBasePage(user, nosurf.Token(r)),
	}

	coursesList, author, err := h.CreateCoursesList(r)
	if err != nil {
		if err == ErrAuthorNotFound || err == ErrCoursesNotFound {
			if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusNotFound); err != nil {
				h.ErrorLog.Println(err)
			}
		} else {
			h.ErrorLog.Printf("Failed to create courses list: %s\n", err)

			if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}
		}

		return
	}

	pageData.Author = author
	pageData.Courses = coursesList

	lenCourses, err := h.Database.CountCoursesWrittenBy(author.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to count all courses written by \"%s\": %s\n", author.ID, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.LenCourses = lenCourses

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "authors-courses", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CoursesPaginationGet(w http.ResponseWriter, r *http.Request) {
	coursesList, _, err := h.CreateCoursesList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create courses list: %s\n", err)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "courses", html.TutorialsListComponent{ErrorMessage: "Failed to get courses. Please try again."}); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "courses", coursesList); err != nil {
		h.ErrorLog.Println(err)
	}
}

// Possible URL queries:
// -page
// -query
func (h *Handlers) CreateCoursesList(r *http.Request) (*html.CoursesListComponent, *models.UserModel, error) {
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

	courses, err := h.Database.GetCourses(query, author.ID, page, CoursesPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all courses (page %d): %s\n", page, err)
		return nil, nil, err
	}

	if len(courses) == 0 {
		return nil, nil, ErrCoursesNotFound
	}

	var coursesSlice []*models.CourseModel
	var lastCourse *models.CourseModel

	if len(courses) < CoursesPerPagination {
		coursesSlice = courses
	} else {
		coursesSlice = courses[:len(courses)-1]
		lastCourse = courses[len(courses)-1]
	}

	tutorialList := &html.CoursesListComponent{
		Courses:    coursesSlice,
		LastCourse: lastCourse,
		QueryURL:   fmt.Sprintf("/authors/%s/courses/htmx?%s", author.Slug, urlQuery.Encode()),
	}

	return tutorialList, author, nil
}
