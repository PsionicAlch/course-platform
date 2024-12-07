package forms

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/payments"
	"github.com/PsionicAlch/psionicalch-home/website/forms/validators"
	"github.com/PsionicAlch/psionicalch-home/website/html"
)

func ValidateAffiliateCode(user *models.UserModel, payment *payments.Payments) validators.ValidationFunc {
	return func(data string, values url.Values) error {
		if data != "" {
			_, err := payment.ValidateAffiliateCode(user.ID, data)
			if err != nil {
				switch err {
				case payments.ErrCantUseOwnAffiliateCode:
					return errors.New("you can't use your own affiliate code")
				case payments.ErrInvalidAffiliateCode:
					return errors.New("invalid affiliate code")
				default:
					return errors.New("failed to validate affiliate code")
				}
			}
		}

		return nil
	}
}

func ValidateAffiliatePoints(user *models.UserModel, payment *payments.Payments) validators.ValidationFunc {
	return func(data string, values url.Values) error {
		if data != "" {
			pointsUsed, err := strconv.ParseUint(data, 10, 64)
			if err != nil {
				return errors.New("only numbers are allowed")
			}

			_, err = payment.ValidateAffiliatePointsUsed(user.ID, uint(pointsUsed))
			if err != nil {
				switch err {
				case payments.ErrInsufficientAffiliatePoints:
					return errors.New("you do not have enough affiliate points available")
				default:
					return errors.New("failed to validate affiliate points")
				}
			}
		}

		return nil
	}
}

func ValidateDiscountCode(payment *payments.Payments) validators.ValidationFunc {
	return func(data string, values url.Values) error {
		if data != "" {
			_, err := payment.ValidateDiscountCode(data)
			if err != nil {
				switch err {
				case payments.ErrInvalidDiscountCode:
					return errors.New("invalid discount code")
				default:
					return errors.New("failed to validate discount code")
				}
			}
		}

		return nil
	}
}

func NewCoursePurchaseForm(r *http.Request, user *models.UserModel, payment *payments.Payments) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		AffiliateCodeName:   ValidateAffiliateCode(user, payment),
		AffiliatePointsName: ValidateAffiliatePoints(user, payment),
		DiscountCodeName:    ValidateDiscountCode(payment),
	})
}

func EmptyCoursePurchaseFormComponent(courseSlug string, user *models.UserModel) *html.CoursePurchaseFormComponent {
	validationURL := fmt.Sprintf("/courses/%s/purchase/validate", courseSlug)

	affiliateCodeInput := new(html.FormControlComponent)
	affiliateCodeInput.Label = "Affiliate Code:"
	affiliateCodeInput.Type = "text"
	affiliateCodeInput.Name = AffiliateCodeName
	affiliateCodeInput.ValidationURL = validationURL

	affiliatePointsInput := new(html.FormControlComponent)
	affiliatePointsInput.Label = fmt.Sprintf("Affiliate Points (%d points available):", user.AffiliatePoints)
	affiliatePointsInput.Type = "number"
	affiliatePointsInput.Name = AffiliatePointsName
	affiliatePointsInput.ValidationURL = validationURL

	discountCodeInput := new(html.FormControlComponent)
	discountCodeInput.Label = "Discount Code:"
	discountCodeInput.Type = "text"
	discountCodeInput.Name = DiscountCodeName
	discountCodeInput.ValidationURL = validationURL

	coursePurchaseForm := new(html.CoursePurchaseFormComponent)
	coursePurchaseForm.AffiliateCodeInput = affiliateCodeInput
	coursePurchaseForm.AffiliatePointsInput = affiliatePointsInput
	coursePurchaseForm.DiscountCodeInput = discountCodeInput
	coursePurchaseForm.CourseSlug = courseSlug
	coursePurchaseForm.Total = payments.CoursePrice

	return coursePurchaseForm
}

func NewCoursePurchaseFormComponent(form *GenericForm, courseSlug string, user *models.UserModel, payment *payments.Payments) *html.CoursePurchaseFormComponent {
	coursePurchaseForm := EmptyCoursePurchaseFormComponent(courseSlug, user)

	affiliateCode := form.GetValue(AffiliateCodeName)
	affiliateCodeErrors := form.GetErrors(AffiliateCodeName)

	coursePurchaseForm.AffiliateCodeInput.Value = affiliateCode
	coursePurchaseForm.AffiliateCodeInput.Errors = affiliateCodeErrors

	if len(affiliateCodeErrors) == 0 {
		if affiliateCode != "" {
			if affiliateCodeDiscount, err := payment.ValidateAffiliateCode(user.ID, affiliateCode); err == nil {
				coursePurchaseForm.AffiliateCodeDiscount = uint(affiliateCodeDiscount * 100)
			}
		}
	}

	affiliatePoints := form.GetValue(AffiliatePointsName)
	affiliatePointsErrors := form.GetErrors(AffiliatePointsName)

	coursePurchaseForm.AffiliatePointsInput.Value = affiliatePoints
	coursePurchaseForm.AffiliatePointsInput.Errors = affiliatePointsErrors

	if len(affiliatePointsErrors) == 0 {
		if affiliatePoints != "" {
			if aps, err := strconv.ParseUint(affiliatePoints, 10, 64); err == nil {
				if affiliatePointsDiscount, err := payment.ValidateAffiliatePointsUsed(user.ID, uint(aps)); err == nil {
					coursePurchaseForm.AffiliatePointsDiscount = uint(affiliatePointsDiscount * 100)
				}
			}
		}
	}

	discountCode := form.GetValue(DiscountCodeName)
	discountCodeErrors := form.GetErrors(DiscountCodeName)

	coursePurchaseForm.DiscountCodeInput.Value = discountCode
	coursePurchaseForm.DiscountCodeInput.Errors = discountCodeErrors

	if len(discountCodeErrors) == 0 {
		if discountCode != "" {
			if discountCodeDiscount, err := payment.ValidateDiscountCode(discountCode); err == nil {
				coursePurchaseForm.DiscountCodeDiscount = uint(discountCodeDiscount * 100)
			}
		}
	}

	if aps, err := strconv.ParseUint(affiliatePoints, 10, 64); err == nil {
		if total, err := payment.CalculatePrice(user.ID, affiliateCode, discountCode, uint(aps)); err == nil {
			coursePurchaseForm.Total = float64(total) / 100.0
		}
	} else {
		if total, err := payment.CalculatePrice(user.ID, affiliateCode, discountCode, 0); err == nil {
			coursePurchaseForm.Total = float64(total) / 100.0
		}
	}

	return coursePurchaseForm
}

func GetCoursePurchaseFormValues(form *GenericForm) (affiliateCode, discountCode string, affiliatePointsUsed uint) {
	affiliateCode = form.GetValue(AffiliateCodeName)
	discountCode = form.GetValue(DiscountCodeName)

	if aps, err := strconv.ParseUint(form.GetValue(AffiliatePointsName), 10, 64); err == nil {
		affiliatePointsUsed = uint(aps)
	} else {
		affiliatePointsUsed = 0
	}

	return
}
