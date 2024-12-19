package sitemap

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/pkg/sitemapper"
)

type Handlers struct {
	Mapper *sitemapper.SiteMapper
}

func SetupHandlers(mapper *sitemapper.SiteMapper) *Handlers {
	return &Handlers{
		Mapper: mapper,
	}
}

func (h *Handlers) SitemapGet(w http.ResponseWriter, r *http.Request) {
	sitemap := h.Mapper.GenerateSitemap("https://www.psionicalch.com")

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(sitemap))
}
