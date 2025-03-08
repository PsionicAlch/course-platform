package rss

import (
	"net/http"

	"github.com/PsionicAlch/course-platform/web/pages"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	*pages.HandlerContext
}

func SetupHandlers(handlerContext *pages.HandlerContext) *Handlers {
	return &Handlers{
		HandlerContext: handlerContext,
	}
}

func (h *Handlers) RSSGet(w http.ResponseWriter, r *http.Request) {
	feed := h.Cache.GetGeneralRSSFeed()

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(feed))
}

func (h *Handlers) RSSTutorialsGet(w http.ResponseWriter, r *http.Request) {
	feed := h.Cache.GetTutorialsRSSFeed()

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(feed))
}

func (h *Handlers) RSSTutorialGet(w http.ResponseWriter, r *http.Request) {
	tutorialSlug := chi.URLParam(r, "tutorial-slug")
	feed := h.Cache.GetTutorialRSSFeed(tutorialSlug)

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(feed))
}

func (h *Handlers) RSSCoursesGet(w http.ResponseWriter, r *http.Request) {
	feed := h.Cache.GetCoursesRSSFeed()

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(feed))
}

func (h *Handlers) RSSTutorialAuthorGet(w http.ResponseWriter, r *http.Request) {
	authorSlug := chi.URLParam(r, "author-slug")
	feed := h.Cache.GetAuthorTutorialsRSSFeed(authorSlug)

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(feed))
}

func (h *Handlers) RSSCourseAuthorGet(w http.ResponseWriter, r *http.Request) {
	authorSlug := chi.URLParam(r, "author-slug")
	feed := h.Cache.GetAuthorCoursesRSSFeed(authorSlug)

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(feed))
}
