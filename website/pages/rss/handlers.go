package rss

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/cache"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Cache cache.Cache
}

func SetupHandlers(c cache.Cache) *Handlers {
	return &Handlers{
		Cache: c,
	}
}

func (h *Handlers) RSSGet(w http.ResponseWriter, r *http.Request) {
	feed := h.Cache.GetRSSFeed()

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(feed))
}

func (h *Handlers) RSSTutorialsGet(w http.ResponseWriter, r *http.Request) {
	feed := h.Cache.GetTutorialsRSSFeed()

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
