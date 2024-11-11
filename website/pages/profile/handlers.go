package profile

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/pkg/gatekeeper"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
)

type Handlers struct {
	utils.Loggers
	renderers pages.Renderers
	auth      *gatekeeper.Gatekeeper
	db        database.Database
}

func SetupHandlers(pageRenderer render.Renderer, auth *gatekeeper.Gatekeeper, db database.Database) *Handlers {
	loggers := utils.CreateLoggers("PROFILE HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		renderers: *pages.CreateRenderers(pageRenderer, nil),
		auth:      auth,
		db:        db,
	}
}

func (h *Handlers) ProfileGet(w http.ResponseWriter, r *http.Request) {
	h.renderers.Page.RenderHTML(w, "profile.page.tmpl", nil)
}

func (h *Handlers) AffiliateHistoryGet(w http.ResponseWriter, r *http.Request) {
	h.renderers.Page.RenderHTML(w, "affiliate-history.page.tmpl", nil)
}

func (h *Handlers) CoursesGet(w http.ResponseWriter, r *http.Request) {
	h.renderers.Page.RenderHTML(w, "profile-courses.page.tmpl", nil)
}

func (h *Handlers) CourseGet(w http.ResponseWriter, r *http.Request) {
	// Redirect to the first incomplete chapter or the last chapter if all are complete.
	utils.Redirect(w, r, "/profile/courses/course-slug-goes-here/chapter-slug-goes-here")
}

func (h *Handlers) CourseChapterGet(w http.ResponseWriter, r *http.Request) {
	// Render the current chapter based off the course slug and chapter slug.
	h.renderers.Page.RenderHTML(w, "profile-course.page.tmpl", nil)
}

func (h *Handlers) TutorialsBookmarksGet(w http.ResponseWriter, r *http.Request) {
	h.renderers.Page.RenderHTML(w, "profile-tutorials-bookmarked.page.tmpl", nil)
}

func (h *Handlers) TutorialsLikedGet(w http.ResponseWriter, r *http.Request) {
	h.renderers.Page.RenderHTML(w, "profile-tutorials-liked.page.tmpl", nil)
}
