package authors

import (
	"fmt"
	"net/http"

	"github.com/PsionicAlch/course-platform/internal/utils"
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

func (h *Handlers) AuthorGet(w http.ResponseWriter, r *http.Request) {
	authorSlug := chi.URLParam(r, "author-slug")
	utils.Redirect(w, r, fmt.Sprintf("/authors/%s/tutorials", authorSlug))
}
