package profile

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	renderers pages.Renderers
	db        database.Database
}

func SetupHandlers(pageRenderer render.Renderer, db database.Database) *Handlers {
	loggers := utils.CreateLoggers("PROFILE HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		renderers: *pages.CreateRenderers(pageRenderer, nil),
		db:        db,
	}
}

func (h *Handlers) ProfileGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "profile.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) AffiliateHistoryGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "affiliate-history.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CoursesGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "profile-courses.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CourseGet(w http.ResponseWriter, r *http.Request) {
	// Redirect to the first incomplete chapter or the last chapter if all are complete.
	utils.Redirect(w, r, "/profile/courses/course-slug-goes-here/chapter-slug-goes-here")
}

func (h *Handlers) CourseChapterGet(w http.ResponseWriter, r *http.Request) {
	// Render the current chapter based off the course slug and chapter slug.
	err := h.renderers.Page.RenderHTML(w, "profile-course.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialsBookmarksGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "profile-tutorials-bookmarks.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) TutorialsLikedGet(w http.ResponseWriter, r *http.Request) {
	err := h.renderers.Page.RenderHTML(w, "profile-tutorials-liked.page.tmpl", nil)
	if err != nil {
		h.ErrorLog.Println(err)
	}
}
