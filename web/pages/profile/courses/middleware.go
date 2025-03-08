package courses

import (
	"fmt"
	"net/http"

	"github.com/PsionicAlch/course-platform/internal/authentication"
	"github.com/PsionicAlch/course-platform/internal/utils"
	"github.com/PsionicAlch/course-platform/web/html"
	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
)

func (h *Handlers) UserBoughtCourse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := authentication.GetUserFromRequest(r)

		courseSlug := chi.URLParam(r, "course-slug")

		course, err := h.Database.GetCourseBySlug(courseSlug)
		if err != nil {
			h.ErrorLog.Printf("Failed to get course by slug (\"%s\"): %s\n", courseSlug, err)

			if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		if course == nil {
			if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusNotFound); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		if !course.Published {
			if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusNotFound); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		hasPurchased, err := h.Database.HasUserPurchasedCourse(user.ID, course.ID)
		if err != nil {
			h.ErrorLog.Printf("Failed to check if user (\"%s\") has purchased the course (\"%s\"): %s\n", user.ID, course.Title, err)

			if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		if !hasPurchased {
			utils.Redirect(w, r, fmt.Sprintf("/courses/%s/purchase", course.Slug))
			return
		}

		ctx := NewContextWithCourseModel(course, r.Context())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
