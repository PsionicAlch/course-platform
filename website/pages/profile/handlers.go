package profile

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/session"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/html"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

const TutorialsPerPagination = 4

type Handlers struct {
	utils.Loggers
	Renderers pages.Renderers
	Auth      *authentication.Authentication
	Database  database.Database
	Session   *session.Session
}

func SetupHandlers(pageRenderer render.Renderer, htmxRenderer render.Renderer, auth *authentication.Authentication, db database.Database, sessions *session.Session) *Handlers {
	loggers := utils.CreateLoggers("PROFILE HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: *pages.CreateRenderers(pageRenderer, htmxRenderer, nil),
		Auth:      auth,
		Database:  db,
		Session:   sessions,
	}
}

func (h *Handlers) ProfileGet(w http.ResponseWriter, r *http.Request) {
	const Elements = 4

	user := authentication.GetUserFromRequest(r)
	pageData := html.ProfilePage{
		BasePage: html.NewBasePage(user),
		User:     user,
	}

	affiliateHistory, err := h.Database.CountUserAffiliateHistory(user.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to get affiliate points history associated with user (\"%s\"): %s\n", user.ID, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user)}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.HasAffiliateHistory = affiliateHistory != 0

	courses, err := h.Database.GetCoursesBoughtByUser("", user.ID, 1, Elements+2)
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

	tutorialsBookmarked, err := h.Database.GetTutorialsBookmarkedByUser("", user.ID, 1, Elements+2)
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

	tutorialsLiked, err := h.Database.GetTutorialsLikedByUser("", user.ID, 1, Elements+2)
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
