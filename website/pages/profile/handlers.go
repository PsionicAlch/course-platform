package profile

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	Renderers pages.Renderers
	Auth      *authentication.Authentication
	Database  database.Database
}

func SetupHandlers(pageRenderer render.Renderer, auth *authentication.Authentication, db database.Database) *Handlers {
	loggers := utils.CreateLoggers("PROFILE HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: *pages.CreateRenderers(pageRenderer, nil),
		Auth:      auth,
		Database:  db,
	}
}

func (h *Handlers) ProfileGet(w http.ResponseWriter, r *http.Request) {
	const Elements = 4

	user := authentication.GetUserFromRequest(r)
	pageData := html.ProfilePage{
		BasePage: html.NewBasePage(user),
		User:     user,
		// TODO: Determine this based on affiliate point transaction history.
		HasAffiliateHistory: false,
	}

	courses, err := h.Database.GetCoursesBoughtByUser(user.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to get courses bought by user (\"%s\"): %s\n", user.ID, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var coursesSlice []*models.CourseModel
	var hasMoreCourses bool

	if len(courses) <= Elements {
		coursesSlice = courses
		hasMoreCourses = false
	} else {
		coursesSlice = courses[:Elements]
		hasMoreCourses = true
	}

	pageData.Courses = coursesSlice
	pageData.HasMoreCourses = hasMoreCourses

	tutorialsBookmarked, err := h.Database.GetTutorialsBookmarkedByUser(user.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to get tutorials bookmarked by user (\"%s\"): %s\n", user.ID, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var tutorialsBookmarkedSlice []*models.TutorialModel
	var hasMoreTutorialsBookmarked bool

	if len(tutorialsBookmarked) <= Elements {
		tutorialsBookmarkedSlice = tutorialsBookmarked
		hasMoreTutorialsBookmarked = false
	} else {
		tutorialsBookmarkedSlice = tutorialsBookmarked[:Elements]
		hasMoreTutorialsBookmarked = true
	}

	pageData.NumTutorialsBookmarked = uint(len(tutorialsBookmarked))
	pageData.TutorialsBookmarked = tutorialsBookmarkedSlice
	pageData.HasMoreTutorialsBookmarked = hasMoreTutorialsBookmarked

	tutorialsLiked, err := h.Database.GetTutorialsLikedByUser(user.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to get tutorials liked by user (\"%s\"): %s\n", user.ID, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var tutorialsLikedSlice []*models.TutorialModel
	var hasMoreTutorialsLiked bool

	if len(tutorialsLiked) <= Elements {
		tutorialsLikedSlice = tutorialsLiked
		hasMoreTutorialsLiked = false
	} else {
		tutorialsLikedSlice = tutorialsLiked[:Elements]
		hasMoreTutorialsLiked = true
	}

	pageData.TutorialsLiked = tutorialsLikedSlice
	pageData.HasMoreTutorialsLiked = hasMoreTutorialsLiked

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "profile", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) AffiliateHistoryGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.ProfileAffiliateHistoryPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "profile-affiliate-history", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CoursesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.ProfileCourses{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "profile-courses", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CourseGet(w http.ResponseWriter, r *http.Request) {
	// Redirect to the first incomplete chapter or the last chapter if all are complete.
	utils.Redirect(w, r, "/profile/courses/course-slug-goes-here/chapter-slug-goes-here")
}

func (h *Handlers) CourseChapterGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.ProfileCourse{
		BasePage: html.NewBasePage(user),
	}

	// Render the current chapter based off the course slug and chapter slug.
	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "profile-course", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialsBookmarksGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.ProfileTutorialsBookmarksPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "profile-tutorials-bookmarks", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialsLikedGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.ProfileTutorialsLikedPage{
		BasePage: html.NewBasePage(user),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "profile-tutorials-liked", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}
