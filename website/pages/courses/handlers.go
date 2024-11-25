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
	"github.com/PsionicAlch/psionicalch-home/internal/session"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
	"github.com/go-chi/chi/v5"
)

const CoursesPerPagination = 5

type Handlers struct {
	utils.Loggers
	Renderers *pages.Renderers
	Database  database.Database
	Session   *session.Session
	Auth      *authentication.Authentication
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, db database.Database, sessions *session.Session, auth *authentication.Authentication) *Handlers {
	loggers := utils.CreateLoggers("COURSE HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: pages.CreateRenderers(pageRenderer, htmxRenderer),
		Database:  db,
		Session:   sessions,
		Auth:      auth,
	}
}

func (h *Handlers) CoursesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.CoursesPage{
		BasePage: html.NewBasePage(user),
	}

	courses, err := h.Database.GetAllCoursesPaginated(1, CoursesPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all courses (page 1) from the database: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var courseSlice []*models.CourseModel
	var lastCourse *models.CourseModel

	if len(courses) < CoursesPerPagination {
		courseSlice = courses
	} else {
		courseSlice = courses[:len(courses)-1]
		lastCourse = courses[len(courses)-1]
	}

	pageData.Courses = &html.CoursesListComponent{
		Courses:    courseSlice,
		LastCourse: lastCourse,
		QueryURL:   fmt.Sprintf("/courses/page/%d", 2),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "courses", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CoursesPaginationGet(w http.ResponseWriter, r *http.Request) {
	coursesComponent := &html.CoursesListComponent{}

	pageNumber, err := strconv.Atoi(chi.URLParam(r, "page-number"))
	if err != nil {
		h.ErrorLog.Printf("Failed to convert page-number to int: %s\n", err)
		coursesComponent.ErrorMessage = "Unexpected server error. Please try again."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "courses", coursesComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	courses, err := h.Database.GetAllCoursesPaginated(pageNumber, CoursesPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all courses (page %d) from the database: %s\n", pageNumber, err)
		coursesComponent.ErrorMessage = "Failed to fetch next courses."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "courses", coursesComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var courseSlice []*models.CourseModel
	var lastCourse *models.CourseModel

	if len(courses) < CoursesPerPagination {
		courseSlice = courses
	} else {
		courseSlice = courses[:len(courses)-1]
		lastCourse = courses[len(courses)-1]
	}

	coursesComponent.Courses = courseSlice
	coursesComponent.LastCourse = lastCourse
	coursesComponent.QueryURL = fmt.Sprintf("/courses/page/%d", pageNumber+1)

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "courses", coursesComponent); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CoursesSearchGet(w http.ResponseWriter, r *http.Request) {
	var page int

	queryPage := r.URL.Query().Get("page")
	query := r.URL.Query().Get("query")
	coursesComponent := &html.CoursesListComponent{}

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

	courses, err := h.Database.SearchCoursesPaginated(query, page, CoursesPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to search for courses (page %d) from the database: %s\n", page, err)
		coursesComponent.ErrorMessage = "Failed to get courses. Please try again."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "courses", coursesComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var courseSlice []*models.CourseModel
	var lastCourse *models.CourseModel

	if len(courses) < CoursesPerPagination {
		courseSlice = courses
	} else {
		courseSlice = courses[:len(courses)-1]
		lastCourse = courses[len(courses)-1]
	}

	coursesComponent.Courses = courseSlice
	coursesComponent.LastCourse = lastCourse
	coursesComponent.QueryURL = fmt.Sprintf("/courses/search?page=%d&query=%s", page+1, url.QueryEscape(query))

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "courses", coursesComponent); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CourseGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.CoursesCoursePage{
		BasePage: html.NewBasePage(user),
	}

	courseSlug := chi.URLParam(r, "slug")

	course, err := h.Database.GetCourseBySlug(courseSlug)
	if err != nil {
		h.ErrorLog.Printf("Failed to get course from the database with slug \"%s\": %s\n", courseSlug, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if course == nil {
		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors500Page{
			BasePage: html.NewBasePage(user),
		}, http.StatusNotFound); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Course = course

	chapters, err := h.Database.CountChapters(course.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to count all the chapters, connected to course \"%s\", in the database: %s\n", course.Title, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Chapters = chapters

	var authorID string
	if course.AuthorID.Valid {
		authorID = course.AuthorID.String
	} else {
		authorID = ""
	}

	// TODO: Implement author fetching logic.
	pageData.Author = &models.AuthorModel{
		ID:      authorID,
		Name:    "Jean-Jacques",
		Surname: "Strydom",
		Slug:    "jean-jacques-strydom",
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "courses-course", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PurchaseCourseGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.CoursesPurchasesPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "courses-purchase", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}
