package courses

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/web/config"
	"github.com/PsionicAlch/psionicalch-home/web/forms"
	"github.com/PsionicAlch/psionicalch-home/web/html"
	"github.com/PsionicAlch/psionicalch-home/web/pages"
	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
)

// TODO: Fix title for course chapters

const CoursesPerPagination = 25

type Handlers struct {
	utils.Loggers
	*pages.HandlerContext
}

func SetupHandlers(handlerContext *pages.HandlerContext) *Handlers {
	loggers := utils.CreateLoggers("COURSE HANDLERS")

	return &Handlers{
		Loggers:        loggers,
		HandlerContext: handlerContext,
	}
}

func (h *Handlers) CoursesGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.CoursesPage{
		BasePage: html.NewBasePage(user, nosurf.Token(r)),
	}

	coursesList, err := h.CreateCoursesList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create courses list: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Courses = coursesList

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "courses", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CoursesPaginationGet(w http.ResponseWriter, r *http.Request) {
	coursesList, err := h.CreateCoursesList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create courses list: %s\n", err)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "courses", html.CoursesListComponent{ErrorMessage: "Failed to get courses. Please try again."}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "courses", coursesList); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) CourseGet(w http.ResponseWriter, r *http.Request) {
	// TODO: Fill the course price in based on the payment.CoursePrice const

	user := authentication.GetUserFromRequest(r)
	pageData := html.CoursesCoursePage{
		BasePage: html.NewBasePage(user, nosurf.Token(r)),
	}

	courseSlug := chi.URLParam(r, "course-slug")

	course, err := h.Database.GetCourseBySlug(courseSlug)
	if err != nil {
		h.ErrorLog.Printf("Failed to get course from the database with slug \"%s\": %s\n", courseSlug, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user, nosurf.Token(r)),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if course == nil {
		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors500Page{
			BasePage: html.NewBasePage(user, nosurf.Token(r)),
		}, http.StatusNotFound); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Course = course

	chapters, err := h.Database.CountChapters(course.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to count all the chapters, connected to course \"%s\", in the database: %s\n", course.Title, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user, nosurf.Token(r)),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Chapters = chapters

	var authorID string
	if course.AuthorID.Valid {
		authorID = course.AuthorID.String
	} else {
		authorID = ""
	}

	author, err := h.Database.GetUserByID(authorID, database.Author)
	if err != nil {
		h.ErrorLog.Printf("Failed to get author by ID \"%s\", in the database: %s\n", authorID, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user, nosurf.Token(r)),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Author = author

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "courses-course", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PurchaseCourseGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.CoursesPurchasesPage{
		BasePage: html.NewBasePage(user, nosurf.Token(r)),
	}

	courseSlug := chi.URLParam(r, "course-slug")

	course, err := h.Database.GetCourseBySlug(courseSlug)
	if err != nil {
		h.ErrorLog.Printf("Failed to get course by slug: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if course == nil {
		h.ErrorLog.Printf("Could not find a course by slug: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusNotFound); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Course = course

	if !course.AuthorID.Valid {
		h.ErrorLog.Printf("Course does not contain a valid author ID: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusNotFound); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	author, err := h.Database.GetUserByID(course.AuthorID.String, database.Author)
	if err != nil {
		h.ErrorLog.Printf("Failed to get course by slug: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if author == nil {
		h.ErrorLog.Printf("Could not find a author by ID: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-404", html.Errors404Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusNotFound); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Author = author

	hasPurchasedCourse, err := h.Database.HasUserPurchasedCourse(user.ID, course.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to check if user (\"%s\") has purchased course (\"%s\"): %s\n", user.ID, course.ID, err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{BasePage: html.NewBasePage(user, nosurf.Token(r))}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if hasPurchasedCourse {
		utils.Redirect(w, r, fmt.Sprintf("/profile/courses/%s", course.Slug))
		return
	}

	purchaseCourseForm := forms.EmptyCoursePurchaseFormComponent(course.Slug, user)
	pageData.CoursePurchaseForm = purchaseCourseForm

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "courses-purchase", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PurchaseCoursePost(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	courseSlug := chi.URLParam(r, "course-slug")
	coursePurchaseForm := forms.NewCoursePurchaseForm(r, user, h.Payment)

	if !coursePurchaseForm.Validate() {
		if err := h.Renderers.Htmx.RenderHTML(w, nil, "course-purchase-form", forms.NewCoursePurchaseFormComponent(coursePurchaseForm, courseSlug, user, h.Payment)); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	affiliateCode, discountCode, affiliatePointsUsed := forms.GetCoursePurchaseFormValues(coursePurchaseForm)

	totalPrice, err := h.Payment.CalculatePrice(user.ID, affiliateCode, discountCode, affiliatePointsUsed)
	if err != nil {
		h.ErrorLog.Printf("Failed to calculate course price: %s\n", err)

		coursePurchaseFormComponent := forms.NewCoursePurchaseFormComponent(coursePurchaseForm, courseSlug, user, h.Payment)
		coursePurchaseFormComponent.ErrorMessage = "Failed to calculate price. Please try again."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "course-purchase-form", coursePurchaseFormComponent, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	course, err := h.Database.GetCourseBySlug(courseSlug)
	if err != nil {
		h.ErrorLog.Printf("Failed to get course from slug: %s\n", err)

		coursePurchaseFormComponent := forms.NewCoursePurchaseFormComponent(coursePurchaseForm, courseSlug, user, h.Payment)
		coursePurchaseFormComponent.ErrorMessage = "Unexpected server error. Please try again."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "course-purchase-form", coursePurchaseFormComponent, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var domainName string

	if config.InDevelopment() {
		domain := config.GetWithoutError[string]("DOMAIN_NAME")
		port := config.GetWithoutError[string]("PORT")
		domainName = fmt.Sprintf("http://%s:%s", domain, port)
	} else {
		domainName = fmt.Sprintf("https://%s", config.GetWithoutError[string]("DOMAIN_NAME"))
	}

	redirectURL, err := h.Payment.BuyCourse(user, course, fmt.Sprintf("%s/courses/%s/purchase/success", domainName, course.Slug), fmt.Sprintf("%s/courses/%s/purchase/cancel", domainName, course.Slug), affiliateCode, discountCode, affiliatePointsUsed, totalPrice)
	if err != nil {
		h.ErrorLog.Printf("Failed to get course from slug: %s\n", err)

		coursePurchaseFormComponent := forms.NewCoursePurchaseFormComponent(coursePurchaseForm, courseSlug, user, h.Payment)
		coursePurchaseFormComponent.ErrorMessage = "Unexpected server error. Please try again."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "course-purchase-form", coursePurchaseFormComponent, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	w.Header().Set("HX-Redirect", redirectURL)
}

func (h *Handlers) PurchaseCourseSuccessGet(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if !h.Payment.ValidatePaymentToken(token) {
		utils.Redirect(w, r, "/courses")
		return
	}

	courseSlug := chi.URLParam(r, "course-slug")

	loadingScreen := &html.LoadingScreenPage{
		Title:   "Validating Purchase",
		PingURL: fmt.Sprintf("/courses/%s/purchase/check?token=%s", courseSlug, token),
	}

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "loading-screen", loadingScreen); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PurchaseCourseCancelGet(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if !h.Payment.ValidatePaymentToken(token) {
		utils.Redirect(w, r, "/courses")
		return
	}

	if err := h.Payment.DeletePaymentToken(token); err != nil {
		h.ErrorLog.Printf("Failed to delete payment token: %s\n", err)
	}

	courseSlug := chi.URLParam(r, "course-slug")

	h.Session.SetWarningMessage(r.Context(), "Payment was cancelled.")

	redirectScreen := &html.RedirectScreenPage{
		RedirectURL: fmt.Sprintf("/courses/%s", courseSlug),
	}

	if err := h.Renderers.Page.RenderHTML(w, nil, "redirect-screen", redirectScreen); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) PurchaseCourseCheckGet(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if !h.Payment.ValidatePaymentToken(token) {
		utils.Redirect(w, r, "/courses")
		return
	}

	user, err := h.Payment.GetUserFromPaymentToken(token)
	if err != nil {
		h.ErrorLog.Printf("Failed to get user from payment token: %s\n", err)
		return
	}

	courseSlug := chi.URLParam(r, "course-slug")

	course, err := h.Database.GetCourseBySlug(courseSlug)
	if err != nil {
		h.ErrorLog.Printf("Failed to get course by slug (\"%s\"): %s\n", courseSlug, err)
		return
	}

	if course == nil {
		utils.Redirect(w, r, "/courses")
		return
	}

	bought, err := h.Database.HasUserPurchasedCourse(user.ID, course.ID)
	if err != nil {
		h.ErrorLog.Printf("Failed to check if the user (\"%s\") has bought this course (\"%s\"): %s\n", user.ID, course.ID, err)
		return
	}

	if bought {
		h.Payment.DeletePaymentToken(token)
		h.Session.SetInfoMessage(r.Context(), "Thank you for your purchase! We hope you enjoy the course.")

		// TODO: Send email thanking user for the purchase of the course.

		utils.Redirect(w, r, fmt.Sprintf("/profile/courses/%s", course.Slug))
	}
}

func (h *Handlers) ValidatePurchasePost(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	courseSlug := chi.URLParam(r, "course-slug")
	coursePurchaseForm := forms.NewCoursePurchaseForm(r, user, h.Payment)
	coursePurchaseForm.Validate()

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "course-purchase-form", forms.NewCoursePurchaseFormComponent(coursePurchaseForm, courseSlug, user, h.Payment)); err != nil {
		h.ErrorLog.Println(err)
	}
}

// Possible URL queries:
// -page
// -query
func (h *Handlers) CreateCoursesList(r *http.Request) (*html.CoursesListComponent, error) {
	query := r.URL.Query().Get("query")
	page := 1

	urlQuery := make(url.Values)

	if pageNum, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		page = pageNum
	}

	urlQuery.Add("query", query)
	urlQuery.Add("page", fmt.Sprintf("%d", page+1))

	courses, err := h.Database.GetCourses(query, "", page, CoursesPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get all courses (page %d) from the database: %s\n", page, err)
		return nil, err
	}

	var coursesSlice []*models.CourseModel
	var lastCourse *models.CourseModel

	if len(courses) < CoursesPerPagination {
		coursesSlice = courses
	} else {
		coursesSlice = courses[:len(courses)-1]
		lastCourse = courses[len(courses)-1]
	}

	coursesList := &html.CoursesListComponent{
		Courses:    coursesSlice,
		LastCourse: lastCourse,
		QueryURL:   fmt.Sprintf("/courses/htmx?%s", urlQuery.Encode()),
	}

	return coursesList, nil
}
