package profile

import (
	"net/http"

	"github.com/PsionicAlch/course-platform/internal/authentication"
	"github.com/PsionicAlch/course-platform/internal/database/models"
	"github.com/PsionicAlch/course-platform/internal/utils"
	"github.com/PsionicAlch/course-platform/web/html"
	"github.com/PsionicAlch/course-platform/web/pages"
	"github.com/justinas/nosurf"
)

const TutorialsPerPagination = 4

type Handlers struct {
	utils.Loggers
	*pages.HandlerContext
}

func SetupHandlers(handlerContext *pages.HandlerContext) *Handlers {
	loggers := utils.CreateLoggers("PROFILE HANDLERS")

	return &Handlers{
		Loggers:        loggers,
		HandlerContext: handlerContext,
	}
}

func (h *Handlers) ProfileGet(w http.ResponseWriter, r *http.Request) {
	const Elements = 4

	user := authentication.GetUserFromRequest(r)
	pageData := html.ProfilePage{
		BasePage: html.NewBasePage(user, nosurf.Token(r)),
		User:     user,
	}

	affiliateHistory, err := h.Database.CountUserAffiliateHistory(user.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to get affiliate points history associated with user (\"%s\"): %s\n", user.ID, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.HasAffiliateHistory = affiliateHistory != 0

	courses, err := h.Database.GetCoursesBoughtByUser("", user.ID, 1, Elements+2)
	if err != nil {
		h.ErrorLog.Printf("Failed to get courses bought by user (\"%s\"): %s\n", user.ID, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
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

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
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

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
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
