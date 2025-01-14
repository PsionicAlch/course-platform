package sitemap

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/web/pages"
)

type Handlers struct {
	*pages.HandlerContext
}

func SetupHandlers(handlerContext *pages.HandlerContext) *Handlers {
	return &Handlers{
		HandlerContext: handlerContext,
	}
}

func (h *Handlers) SitemapGet(w http.ResponseWriter, r *http.Request) {
	sitemap := h.Mapper.GenerateSitemap("https://www.psionicalch.com", "/htmx")

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(sitemap))
}
