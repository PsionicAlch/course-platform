package payments

import "github.com/PsionicAlch/psionicalch-home/internal/database/models"

type Emailer interface {
	SendThankYouForPurchaseEmail(email, firstName, affiliateCode string, course *models.CourseModel, discount *models.DiscountModel)
}
