package sitemap

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/web/pages"
)

type Handlers struct {
	utils.Loggers
	*pages.HandlerContext
}

func SetupHandlers(handlerContext *pages.HandlerContext) *Handlers {
	return &Handlers{
		Loggers:        utils.CreateLoggers("SITEMAP HANDLER"),
		HandlerContext: handlerContext,
	}
}

func (h *Handlers) SitemapGet(w http.ResponseWriter, r *http.Request) {
	sitemap, err := h.Mapper.GenerateSitemap("https://www.psionicalch.com", "/htmx")
	if err != nil {
		h.ErrorLog.Println(err)
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte(sitemap))
}
