package payments

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

type Payments struct {
	utils.Loggers
	WebhookSecret string
	Database      database.Database
}

func SetupPayments(privateKey, webhookSecret string, db database.Database) *Payments {
	loggers := utils.CreateLoggers("PAYMENTS")

	stripe.Key = privateKey

	return &Payments{
		Loggers:       loggers,
		WebhookSecret: webhookSecret,
		Database:      db,
	}
}

func (payment *Payments) BuyCourse(user *models.UserModel, course *models.CourseModel, successUrl, cancelUrl, affiliateCode, discountCode string, affiliatePointsUsed, amountPaid int64) (string, error) {
	paymentKey, err := GeneratePaymentKey()
	if err != nil {
		payment.ErrorLog.Printf("Failed to generate payment key: %s\n", err)
		return "", err
	}

	metaData := map[string]string{
		"user_id":      user.ID,
		"user_name":    user.Name,
		"user_surname": user.Surname,
		"user_email":   user.Email,
		"payment_key":  paymentKey,
	}

	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyUSD)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String(course.Title),
						Description: stripe.String(course.Description),
						// TODO: Make sure that images are hosted with their full URL path.
						// Images: stripe.StringSlice([]string{course.ThumbnailURL}),
					},
					UnitAmount: stripe.Int64(amountPaid),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:    stripe.String(successUrl),
		CancelURL:     stripe.String(cancelUrl),
		CustomerEmail: stripe.String(user.Email),
		Metadata:      metaData,
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: metaData,
		},
	}

	s, err := session.New(params)
	if err != nil {
		payment.ErrorLog.Printf("Failed to create new Stripe checkout session: %s\n", err)
		return "", err
	}

	ac := database.NewNullString(affiliateCode)
	dc := database.NewNullString(discountCode)

	if err := payment.Database.RegisterCoursePurchase(user.ID, course.ID, paymentKey, s.ID, ac, dc, affiliatePointsUsed, float64(amountPaid)/100.0); err != nil {
		payment.ErrorLog.Printf("Failed to save stripe checkout information to the database: %s\n", err)
		return "", err
	}

	return s.URL, nil
}
