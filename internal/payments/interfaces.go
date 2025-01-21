package payments

import "github.com/PsionicAlch/psionicalch-home/internal/database/models"

type Emailer interface {
	SendThankYouForPurchaseEmail(email, firstName, affiliateCode string, course *models.CourseModel, discount *models.DiscountModel)
	SendRefundRequestFailedEmail(email, firstName, courseName, failureReason string)
	SendRefundRequestCancelledEmail(email, firstName, courseName string)
	SendRefundRequestSuccessfulEmail(email, firstName, courseName string, refundAmount float64)
}
