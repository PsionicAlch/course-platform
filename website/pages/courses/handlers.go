package courses

import (
	"fmt"
	"net/http"
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

		h.Session.SetErrorMessage(r.Context(), "Failed to load courses. Please try again")

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user),
		}); err != nil {
			h.ErrorLog.Println(err)
		}
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

func (h *Handlers) CourseGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.CoursesCoursePage{
		BasePage: html.NewBasePage(user),
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
