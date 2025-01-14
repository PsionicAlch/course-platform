package discounts

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/web/forms"
	"github.com/PsionicAlch/psionicalch-home/web/html"
	"github.com/PsionicAlch/psionicalch-home/web/pages"
	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
)

const DiscountPerPagination = 25

var DiscountStatuses = []string{"Active", "Inactive"}

type Handlers struct {
	utils.Loggers
	*pages.HandlerContext
}

// TODO: Set "Used" based on database
func SetupHandlers(handlerContext *pages.HandlerContext) *Handlers {
	loggers := utils.CreateLoggers("ADMIN DISCOUNTS HANDLERS")

	return &Handlers{
		Loggers:        loggers,
		HandlerContext: handlerContext,
	}
}

func (h *Handlers) DiscountsGet(w http.ResponseWriter, r *http.Request) {
	user := authentication.GetUserFromRequest(r)
	pageData := html.AdminDiscountsPage{
		BasePage: html.NewBasePage(user, nosurf.Token(r)),
	}

	numDiscounts, err := h.Database.CountDiscounts()
	if err != nil {
		h.ErrorLog.Printf("Failed to count the number of discounts in the database: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user, nosurf.Token(r)),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.NumDiscounts = numDiscounts
	pageData.DiscountStatus = DiscountStatuses

	discounts, urlQuery, err := h.CreateDiscountsList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create discounts list: %s\n", err)

		if err := h.Renderers.Page.RenderHTML(w, r.Context(), "errors-500", html.Errors500Page{
			BasePage: html.NewBasePage(user, nosurf.Token(r)),
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	pageData.Discounts = discounts

	urlQuery.Set("page", "1")
	pageData.URLQuery = urlQuery.Encode()

	pageData.NewDiscountForm = forms.EmptyNewDiscountFormComponent()

	if err := h.Renderers.Page.RenderHTML(w, r.Context(), "admin-discounts", pageData); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) DiscountsPaginationGet(w http.ResponseWriter, r *http.Request) {
	discounts, _, err := h.CreateDiscountsList(r)
	if err != nil {
		h.ErrorLog.Printf("Failed to create discounts list: %s\n", err)

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "admin-discounts", html.AdminDiscountsListComponent{
			ErrorMessage: "Unexpected server error. Failed to get discounts.",
		}, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "admin-discounts", discounts); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) NewDiscountPost(w http.ResponseWriter, r *http.Request) {
	form := forms.NewDiscountForm(r)

	if !form.Validate() {
		if err := h.Renderers.Htmx.RenderHTML(w, nil, "new-discount-form", forms.NewDiscountFormComponent(form)); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	title, description, uses, amount := forms.GetNewDiscountFormValues(form)
	if _, err := h.Database.AddDiscount(title, description, amount, uses); err != nil {
		formComponent := forms.NewDiscountFormComponent(form)
		formComponent.ErrorMessage = "Unexpected server error. Failed to create new discount."

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "new-discount-form", formComponent); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	w.Header().Set("HX-Redirect", "/admin/discounts")

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "new-discount-form", forms.EmptyNewDiscountFormComponent()); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) ValidateNewDiscountPost(w http.ResponseWriter, r *http.Request) {
	form := forms.NewDiscountFormPartialValidation(r)
	form.Validate()

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "new-discount-form", forms.NewDiscountFormComponent(form)); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) EmptyNewDiscountGet(w http.ResponseWriter, r *http.Request) {
	form := forms.EmptyNewDiscountFormComponent()

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "new-discount-form", form); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) StatusEditGet(w http.ResponseWriter, r *http.Request) {
	discountId := chi.URLParam(r, "discount-id")

	discount, err := h.Database.GetDiscountByID(discountId)
	if err != nil {
		h.ErrorLog.Printf("Failed to get discount \"%s\" from the database: %s\n", discountId, err)

		resp := "Inactive"
		resp += `
		<script>
			notyf.open({
				type: 'flash-error',
				message: "Unexpected server error"
			});
		</script>
		`

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if discount == nil {
		h.ErrorLog.Printf("Failed to get discount by ID \"%s\": Nill was returned\n", discountId)

		resp := "Inactive"
		resp += `
		<script>
			notyf.open({
				type: 'flash-error',
				message: "Discount doesn't exist."
			});
		</script>
		`

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	var discountStatus string
	if discount.Active {
		discountStatus = "Active"
	} else {
		discountStatus = "Inactive"
	}

	discountStatuses := make(map[string]string, len(DiscountStatuses))
	for _, status := range DiscountStatuses {
		discountStatuses[status] = status
	}

	selectComponent := &html.SelectComponent{
		Name:     "discount-status",
		Options:  discountStatuses,
		Selected: discountStatus,
		URL:      fmt.Sprintf("/admin/discounts/htmx/change-status/%s", discountId),
	}

	if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", selectComponent); err != nil {
		h.ErrorLog.Println(err)
	}
}

func (h *Handlers) StatusEditPost(w http.ResponseWriter, r *http.Request) {
	discountId := chi.URLParam(r, "discount-id")

	r.ParseForm()
	discountStatus := r.Form.Get("discount-status")

	discount, err := h.Database.GetDiscountByID(discountId)
	if err != nil {
		h.ErrorLog.Printf("Failed to get discount \"%s\" from the database: %s\n", discountId, err)

		resp := "Inactive"
		resp += `
		<script>
			notyf.open({
				type: 'flash-error',
				message: "Unexpected server error"
			});
		</script>
		`

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if discount == nil {
		h.ErrorLog.Printf("Failed to get discount by ID \"%s\": Nill was returned\n", discountId)

		resp := "Inactive"
		resp += `
		<script>
			notyf.open({
				type: 'flash-error',
				message: "Discount doesn't exist."
			});
		</script>
		`

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", resp, http.StatusInternalServerError); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if !utils.InSlice(discountStatus, DiscountStatuses) {
		discountStatuses := make(map[string]string, len(DiscountStatuses))
		for _, status := range DiscountStatuses {
			discountStatuses[status] = status
		}

		selectComponent := &html.SelectComponent{
			Name:    "discount-status",
			Options: discountStatuses,
			URL:     fmt.Sprintf("/admin/discounts/htmx/change-status/%s", discountId),
		}

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", selectComponent, http.StatusBadRequest); err != nil {
			h.ErrorLog.Println(err)
		}

		return
	}

	if discountStatus == "Active" {
		if err := h.Database.ActivateDiscount(discountId); err != nil {
			h.ErrorLog.Printf("Failed to update discount active status: %s\n", err)

			discountStatuses := make(map[string]string, len(DiscountStatuses))
			for _, status := range DiscountStatuses {
				discountStatuses[status] = status
			}

			selectComponent := &html.SelectComponent{
				Name:         "discount-status",
				Options:      discountStatuses,
				URL:          fmt.Sprintf("/admin/discounts/htmx/change-status/%s", discountId),
				ErrorMessage: "Unexpected server error. Failed to update discount active status.",
			}

			if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", selectComponent, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", "Active"); err != nil {
			h.ErrorLog.Println(err)
		}
	} else {
		if err := h.Database.DeactivateDiscount(discountId); err != nil {
			h.ErrorLog.Printf("Failed to update discount active status: %s\n", err)

			discountStatuses := make(map[string]string, len(DiscountStatuses))
			for _, status := range DiscountStatuses {
				discountStatuses[status] = status
			}

			selectComponent := &html.SelectComponent{
				Name:         "discount-status",
				Options:      discountStatuses,
				URL:          fmt.Sprintf("/admin/discounts/htmx/change-status/%s", discountId),
				ErrorMessage: "Unexpected server error. Failed to update discount active status.",
			}

			if err := h.Renderers.Htmx.RenderHTML(w, nil, "select", selectComponent, http.StatusInternalServerError); err != nil {
				h.ErrorLog.Println(err)
			}

			return
		}

		if err := h.Renderers.Htmx.RenderHTML(w, nil, "empty", "Inactive"); err != nil {
			h.ErrorLog.Println(err)
		}
	}
}

// Possible URL queries:
// -page
// -query
// -status
func (h *Handlers) CreateDiscountsList(r *http.Request) (*html.AdminDiscountsListComponent, url.Values, error) {
	var query string
	var status *bool
	var page int

	urlQuery := make(url.Values)

	if q := r.URL.Query().Get("query"); q != "" {
		query = q

		urlQuery.Add("query", q)
	}

	if s := r.URL.Query().Get("status"); s != "" {
		if s == "Active" {
			temp := true
			status = &temp
		} else if s == "Inactive" {
			temp := false
			status = &temp
		}

		urlQuery.Add("status", s)
	} else {
		status = nil
	}

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil {
		page = p
	} else {
		page = 1
	}

	urlQuery.Add("page", fmt.Sprintf("%d", page+1))

	discounts, err := h.Database.GetDiscountsPaginated(query, status, uint(page), DiscountPerPagination)
	if err != nil {
		h.ErrorLog.Printf("Failed to get discounts from the database: %s\n", err)
		return nil, urlQuery, err
	}

	var discountsSlice []*models.DiscountModel
	var lastDiscount *models.DiscountModel

	if len(discounts) < DiscountPerPagination {
		discountsSlice = discounts
	} else {
		discountsSlice = discounts[:len(discounts)-1]
		lastDiscount = discounts[len(discounts)-1]
	}

	discountUsed := make(map[string]uint, len(discounts))
	for _, discount := range discounts {
		// TODO: Query database for every time this discount code was used.
		discountUsed[discount.ID] = 0
	}

	discountsList := &html.AdminDiscountsListComponent{
		Discounts:    discountsSlice,
		LastDiscount: lastDiscount,
		DiscountUsed: discountUsed,
		BaseURL:      "/admin/discounts/htmx",
		URLQuery:     urlQuery.Encode(),
	}

	return discountsList, urlQuery, nil
}
