package rss

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/cache"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/pages"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	utils.Loggers
	Renderers *pages.Renderers
	Database  database.Database
	Cache     cache.Cache
}

func SetupHandlers(rssRenderer render.Renderer, db database.Database, c cache.Cache) *Handlers {
	loggers := utils.CreateLoggers("RSS HANDLERS")

	return &Handlers{
		Loggers:   loggers,
		Renderers: pages.CreateRenderers(nil, nil, rssRenderer),
		Database:  db,
		Cache:     c,
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
