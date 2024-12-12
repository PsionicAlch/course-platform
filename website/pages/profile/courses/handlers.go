package courses

import (
	"fmt"
	"maps"
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

const CoursesPerPagination = 25

type Handlers struct {
	utils.Loggers
	Renderers pages.Renderers
	Database  database.Database
	Session   *session.Session
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, db database.Database, sessions *session.Session) *Handlers {
	loggers := utils.CreateLoggers("PROFILE HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: *pages.CreateRenderers(pageRenderer, htmxRenderer),
		Database:  db,
		Session:   sessions,
	}
}

func (h *Handlers) CoursesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.ProfileCourses{
		BasePage: html.NewBasePage(user),
	}

	courses, err := h.CreateCoursesList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create courses list: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Courses = courses

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "profile-courses", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CoursesPaginationGet(w http.ResponseWriter, r *http.Request) {
	courses, err := h.CreateCoursesList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create courses list: %s\n", err)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "profile-courses", html.CoursesListComponent{ErrorMessage: "Failed to get courses. Please try again."}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "profile-courses", courses); err != nil {
		h.ErrorLog.Println(err)
	}
}

// Possible URL Queries:
// - page
// - query
func (h *Handlers) CreateCoursesList(r *http.Request) (*html.CoursesListComponent, error) {
	user := authentication.GetUserFromRequest(r)

	page := 1
	query := ""

	urlQuery := make(url.Values)

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		page = p
	}

	urlQuery.Add("page", fmt.Sprintf("%d", page+1))

	if q := r.URL.Query().Get("query"); q != "" {
		query = q

		urlQuery.Add("query", q)
	}

	courses, err := h.Database.GetCoursesBoughtByUser(query, user.ID, uint(page), CoursesPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get courses bought by user (\"%s\"): %s\n", user.ID, err)
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
		QueryURL:   fmt.Sprintf("/profile/courses/htmx?%s", urlQuery.Encode()),
	}

	return coursesList, nil
}

func (h *Handlers) CourseGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	course := GetCourseFromRequest(r)

	if user == nil {
		h.ErrorLog.Println("Failed to get user from context")

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if course == nil {
		h.ErrorLog.Println("Failed to get course from context")

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	chapters, err := h.Database.GetCourseChapters(course.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all course (\"%s\") chapters: %s\n", course.Title, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	for _, chapter := range chapters {
		hasCompleted, err := h.Database.HasUserCompletedChapter(user.ID, course.ID, chapter.ID)
		if err != nil {
			h.ErrorLog.Printf("Could not check if user has completed course chapter: %s\n", err)

			if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		if !hasCompleted {
			utils.Redirect(w, r, fmt.Sprintf("/profile/courses/%s/%s", course.Slug, chapter.Slug))
			return
		}
	}

	// Redirect to the first incomplete chapter or the last chapter if all are complete.
	utils.Redirect(w, r, fmt.Sprintf("/profile/courses/%s/certificate", course.Slug))
}

func (h *Handlers) CourseCertificateGet(w http.ResponseWriter, r *http.Request) {
	// TODO: Figure out how the to generate certificate for users.
}

func (h *Handlers) CourseChapterGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	course := GetCourseFromRequest(r)

	if user == nil {
		h.ErrorLog.Println("Failed to get user from context")

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if course == nil {
		h.ErrorLog.Println("Failed to get course from context")

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData := html.ProfileCourse{
		BasePage: html.NewBasePage(user),
		Course:   course,
	}

	chapterSlug := chi.URLParam(r, "chapter-slug")

	chapter, err := h.Database.GetChapterBySlug(chapterSlug)
	if err != nil {
		h.ErrorLog.Printf("Failed to get chapter by slug (\"%s\"): %s\n", chapterSlug, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if chapter == nil {
		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user)}, http.StatusNotFound); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Chapter = chapter

	chapters, err := h.Database.GetCourseChapters(course.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to get chapters associates with course (\"%s\"): %s\n", course.ID, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if chapters == nil {
		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user)}, http.StatusNotFound); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Chapters = chapters

	if chapters[len(chapters)-1].Slug == chapter.Slug {
		pageData.LastChapter = true
	} else {
		pageData.LastChapter = false
	}

	completed := make(map[string]bool, len(chapters))

	for _, c := range chapters {
		hasCompleted, err := h.Database.HasUserCompletedChapter(user.ID, course.ID, c.ID)
		if err != nil {
			h.ErrorLog.Printf("Failed to check if user (\"%s\") has completed course (\"%s\") chapter (\"%s\"): %s\n", user.ID, course.ID, c.ID, err)

			if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		completed[c.ID] = hasCompleted
	}

	pageData.Completed = completed

	pageData.HasCompletedCourse = !utils.InSeq(false, maps.Values(completed))

	// Render the current chapter based off the course slug and chapter slug.
	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "profile-course", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}
