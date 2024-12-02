package forms

import (
	"math"
	"net/http"
	"strconv"

	"github.com/PsionicAlch/psionicalch-home/website/forms/validators"
	"github.com/PsionicAlch/psionicalch-home/website/html"
)

func NewDiscountForm(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		TitleName: validators.ChainValidators(
			validators.NotEmpty,
			validators.MaxLength(100),
		),
		DescriptionName: validators.ChainValidators(
			validators.NotEmpty,
		),
		UsesName: validators.ChainValidators(
			validators.NotEmpty,
			validators.Integer,
			validators.Min(1),
			validators.Max(math.MaxInt),
		),
		AmountName: validators.ChainValidators(
			validators.NotEmpty,
			validators.Integer,
			validators.Min(1),
			validators.Max(100),
		),
	})
}

func NewDiscountFormPartialValidation(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		TitleName: validators.ChainValidators(
			validators.MaxLength(100),
		),
		DescriptionName: validators.Empty,
		UsesName: validators.ChainValidators(
			validators.Integer,
			validators.Max(math.MaxInt),
		),
		AmountName: validators.ChainValidators(
			validators.NotEmpty,
			validators.Integer,
			validators.Max(100),
		),
	})
}

func EmptyNewDiscountFormComponent() *html.NewDiscountFormComponent {
	titleInput := new(html.FormControlComponent)
	titleInput.Label = "Title:"
	titleInput.Name = TitleName
	titleInput.Type = "text"
	titleInput.ValidationURL = NewDiscountsValidationURL

	descriptionInput := new(html.FormControlComponent)
	descriptionInput.Label = "Description:"
	descriptionInput.Name = DescriptionName
	descriptionInput.Type = "text"
	descriptionInput.ValidationURL = NewDiscountsValidationURL

	usesInput := new(html.FormControlComponent)
	usesInput.Label = "Discount Uses:"
	usesInput.Name = UsesName
	usesInput.Type = "number"
	usesInput.ValidationURL = NewDiscountsValidationURL

	amountInput := new(html.FormControlComponent)
	amountInput.Label = "Discount Amount (%):"
	amountInput.Name = AmountName
	amountInput.Type = "number"
	amountInput.ValidationURL = NewDiscountsValidationURL

	newDiscountForm := new(html.NewDiscountFormComponent)
	newDiscountForm.TitleInput = titleInput
	newDiscountForm.DescriptionInput = descriptionInput
	newDiscountForm.UsesInput = usesInput
	newDiscountForm.AmountInput = amountInput

	return newDiscountForm
}

func NewDiscountFormComponent(form *GenericForm) *html.NewDiscountFormComponent {
	newDiscountForm := EmptyNewDiscountFormComponent()

	newDiscountForm.TitleInput.Value = form.GetValue(TitleName)
	newDiscountForm.TitleInput.Errors = form.GetErrors(TitleName)

	newDiscountForm.DescriptionInput.Value = form.GetValue(DescriptionName)
	newDiscountForm.DescriptionInput.Errors = form.GetErrors(DescriptionName)

	newDiscountForm.UsesInput.Value = form.GetValue(UsesName)
	newDiscountForm.UsesInput.Errors = form.GetErrors(UsesName)

	newDiscountForm.AmountInput.Value = form.GetValue(AmountName)
	newDiscountForm.AmountInput.Errors = form.GetErrors(AmountName)

	return newDiscountForm
}

func GetNewDiscountFormValues(form *GenericForm) (title, description string, uses, amount uint64) {
	title = form.GetValue(TitleName)
	description = form.GetValue(DescriptionName)

	u := form.GetValue(UsesName)
	uses, _ = strconv.ParseUint(u, 10, 64)

	a := form.GetValue(AmountName)
	amount, _ = strconv.ParseUint(a, 10, 64)

	return
}
