package payments

import (
	"fmt"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/refund"
)

type Payments struct {
	utils.Loggers
	WebhookSecret string
	Database      database.Database
	Mailer        Emailer
}

func SetupPayments(privateKey, webhookSecret string, db database.Database, mailer Emailer) *Payments {
	loggers := utils.CreateLoggers("PAYMENTS")

	stripe.Key = privateKey

	return &Payments{
		Loggers:       loggers,
		WebhookSecret: webhookSecret,
		Database:      db,
		Mailer:        mailer,
	}
}

func (payment *Payments) CreateDiscount(title, description string, discountAmount, uses uint64) (*models.DiscountModel, error) {
	discountId, err := payment.Database.AddDiscount(title, description, discountAmount, uses)
	if err != nil {
		payment.ErrorLog.Printf("Failed to create new discount: %s\n", err)
		return nil, err
	}

	if err := payment.Database.ActivateDiscount(discountId); err != nil {
		payment.ErrorLog.Printf("Failed to activate new discount: %s\n", err)
		return nil, err
	}

	discount, err := payment.Database.GetDiscountByID(discountId)
	if err != nil {
		payment.ErrorLog.Printf("Failed to get discount by ID: %s\n", err)
		return nil, err
	}

	return discount, nil
}

func (payment *Payments) BuyCourse(user *models.UserModel, course *models.CourseModel, successUrl, cancelUrl, affiliateCode, discountCode string, affiliatePointsUsed uint, amountPaid int64) (string, error) {
	paymentKey, err := GeneratePaymentKey()
	if err != nil {
		payment.ErrorLog.Printf("Failed to generate payment key: %s\n", err)
		return "", err
	}

	paymentToken, err := database.GenerateToken()
	if err != nil {
		payment.ErrorLog.Printf("Failed to generate payment token: %s\n", err)
		return "", err
	}

	if amountPaid <= 0 {
		ac := database.NewNullString(affiliateCode)
		dc := database.NewNullString(discountCode)

		if err := payment.Database.RegisterCoursePurchase(user.ID, course.ID, paymentKey, "", ac, dc, affiliatePointsUsed, float64(amountPaid)/100.0, "", PaymentToken, time.Now().Add(time.Hour)); err != nil {
			payment.ErrorLog.Printf("Failed register course purchase in the database: %s\n", err)
			return "", err
		}

		coursePurchase, err := payment.Database.GetCoursePurchaseByPaymentKey(paymentKey)
		if err != nil {
			payment.ErrorLog.Printf("Failed to get course purchase by payment key (\"%s\"): %s\n", paymentKey, err)
			return "", err
		}

		if err := payment.Database.UpdateCoursePurchasePaymentStatus(coursePurchase.ID, database.Succeeded); err != nil {
			payment.ErrorLog.Printf("Failed to update course purchase's (\"%s\") payment status: %s\n", coursePurchase.ID, err)
			return "", err
		}

		redirectURL := fmt.Sprintf("/profile/courses/%s", course.Slug)

		if ac.Valid {
			affiliateUser, err := payment.Database.GetUserByAffiliateCode(coursePurchase.AffiliateCode.String, database.All)
			if err != nil {
				payment.ErrorLog.Printf("Failed to get user by affiliate code (\"%s\"): %s\n", coursePurchase.AffiliateCode.String, err)
				return redirectURL, nil
			}

			if err := payment.Database.RegisterAffiliatePointsChange(affiliateUser.ID, coursePurchase.CourseID, AffiliateReward, "Affiliate reward received"); err != nil {
				payment.ErrorLog.Printf("Failed to reward user (\"%s\") with affiliate points: %s\n", affiliateUser.ID, err)
				return redirectURL, nil
			}
		}

		discount, err := payment.CreateDiscount(fmt.Sprintf("Thank You Gift To %s %s", user.Name, user.Surname), "A gift to thank the user for buying a course from us", 20, 1)
		if err != nil {
			payment.ErrorLog.Printf("Failed to create new discount: %s\n", err)
			return redirectURL, nil
		}

		go payment.Mailer.SendThankYouForPurchaseEmail(user.Email, user.Name, user.AffiliateCode, course, discount)

		return redirectURL, nil
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
						Images:      stripe.StringSlice([]string{course.ThumbnailURL}),
					},
					UnitAmount: stripe.Int64(amountPaid),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:    stripe.String(fmt.Sprintf("%s?token=%s", successUrl, paymentToken)),
		CancelURL:     stripe.String(fmt.Sprintf("%s?token=%s", cancelUrl, paymentToken)),
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

	if err := payment.Database.RegisterCoursePurchase(user.ID, course.ID, paymentKey, s.ID, ac, dc, affiliatePointsUsed, float64(amountPaid)/100.0, paymentToken, PaymentToken, time.Now().Add(time.Hour)); err != nil {
		payment.ErrorLog.Printf("Failed to save stripe checkout information to the database: %s\n", err)
		return "", err
	}

	return s.URL, nil
}

func (payment *Payments) RequestRefund(user *models.UserModel, course *models.CourseModel) error {
	coursePurchases, err := payment.Database.GetCoursePurchasesByUserAndCourse(user.ID, course.ID)
	if err != nil {
		payment.ErrorLog.Printf("Failed to get course purchase for user (\"%s\") and course (\"%s\"): %s\n", user.ID, course.ID, err)
		return err
	}

	index, has := utils.Find(coursePurchases, func(coursePurchase *models.CoursePurchaseModel) bool {
		return coursePurchase.PaymentStatus == database.Succeeded.String()
	})

	if !has {
		return ErrUserHasNotBoughtCourse
	}

	coursePurchase := coursePurchases[index]

	if coursePurchase.AmountPaid == 0.0 {
		if err := payment.Database.RegisterRefund(coursePurchase.UserID, coursePurchase.ID, database.RefundSucceeded); err != nil {
			payment.ErrorLog.Printf("Failed to insert new refund: %s\n", err)
			return err
		}

		if err := payment.Database.UpdateCoursePurchasePaymentStatus(coursePurchase.ID, database.Refunded); err != nil {
			payment.ErrorLog.Printf("Failed to update course purchase (\"%s\") status to refunded: %s\n", coursePurchase.ID, err)
			return err
		}

		return nil
	}

	checkoutSession, err := session.Get(coursePurchase.StripeCheckoutSessionID, nil)
	if err != nil {
		payment.ErrorLog.Printf("Failed to get Stripe Checkout Session: %s\n", err)
		return err
	}

	refundParams := &stripe.RefundParams{
		PaymentIntent: stripe.String(checkoutSession.PaymentIntent.ID),
		Metadata: map[string]string{
			"user_id":      user.ID,
			"user_name":    user.Name,
			"user_surname": user.Surname,
			"user_email":   user.Email,
			"payment_key":  coursePurchase.PaymentKey,
		},
	}

	_, err = refund.New(refundParams)
	if err != nil {
		payment.ErrorLog.Printf("Failed to created Stripe Refund: %s\n", err)
		return err
	}

	if err := payment.Database.RegisterRefund(coursePurchase.UserID, coursePurchase.ID, database.RefundPending); err != nil {
		payment.ErrorLog.Printf("Failed to insert new refund: %s\n", err)
		return err
	}

	if err := payment.Database.UpdateCoursePurchasePaymentStatus(coursePurchase.ID, database.Refunded); err != nil {
		payment.ErrorLog.Printf("Failed to update course purchase (\"%s\") status to refunded: %s\n", coursePurchase.ID, err)
	}

	return nil
}
