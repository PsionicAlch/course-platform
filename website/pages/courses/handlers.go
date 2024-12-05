package courses

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/payments"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/session"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
	"github.com/go-chi/chi/v5"
)

// TODO: Clean up!

const CoursesPerPagination = 25

type Handlers struct {
	utils.Loggers
	Renderers *pages.Renderers
	Database  database.Database
	Session   *session.Session
	Auth      *authentication.Authentication
	Payment   *payments.Payments
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, db database.Database, sessions *session.Session, auth *authentication.Authentication, payment *payments.Payments) *Handlers {
	loggers := utils.CreateLoggers("COURSE HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: pages.CreateRenderers(pageRenderer, htmxRenderer),
		Database:  db,
		Session:   sessions,
		Auth:      auth,
		Payment:   payment,
	}
}

func (h *Handlers) CoursesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.CoursesPage{
		BasePage: html.NewBasePage(user),
	}

	coursesList, err := h.CreateCoursesList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create courses list: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Courses = coursesList

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "courses", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CoursesPaginationGet(w http.ResponseWriter, r *http.Request) {
	coursesList, err := h.CreateCoursesList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create courses list: %s\n", err)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "courses", html.CoursesListComponent{ErrorMessage: "Failed to get courses. Please try again."}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "courses", coursesList); err != nil {
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

	author, err := h.Database.GetUserByID(authorID, database.Author)
	if err != nil {
		h.ErrorLog.Printf("Failed to get author by ID \"%s\", in the database: %s\n", authorID, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Author = author

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "courses-course", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PurchaseCourseGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.CoursesPurchasesPage{
		BasePage: html.NewBasePage(user),
	}

	// TODO: Redirect the user to the course chapter if they've already bought this course.

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "courses-purchase", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PurchaseCoursePost(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) PurchaseCourseSuccessGet(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) PurchaseCourseCancelGet(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) ValidatePurchasePost(w http.ResponseWriter, r *http.Request) {

}

// Possible URL queries:
// -page
// -query
func (h *Handlers) CreateCoursesList(r *http.Request) (*html.CoursesListComponent, error) {
	query := r.URL.Query().Get("query")
	page := 1

	urlQuery := make(url.Values)

	if pageNum, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		page = pageNum
	}

	urlQuery.Add("query", query)
	urlQuery.Add("page", fmt.Sprintf("%d", page+1))

	courses, err := h.Database.GetCourses(query, page, CoursesPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all courses (page %d) from the database: %s\n", page, err)
		return nil, err
	}

	var coursesSlice []*models.CourseModel
	var lastCourse *models.CourseModel

	if len(courses) < CoursesPerPagination {
		coursesSlice = courses
	} else {
		coursesSlice = courses[:len(courses)-1]
		lastCourse = courses[len(courses)-1]
	}

	coursesList := &html.CoursesListComponent{
		Courses:    coursesSlice,
		LastCourse: lastCourse,
		QueryURL:   fmt.Sprintf("/courses/htmx?%s", urlQuery.Encode()),
	}

	return coursesList, nil
}
