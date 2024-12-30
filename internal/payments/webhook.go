package payments

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)

func (payment *Payments) Webhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		payment.ErrorLog.Printf("Failed to read request body in webhook: %s\n", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	signatureHeader := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(body, signatureHeader, payment.WebhookSecret)
	if err != nil {
		payment.ErrorLog.Printf("Failed to verify webhook signature: %s\n", err)
		http.Error(w, "invalid signature", http.StatusBadRequest)
		return
	}

	handlers := map[string]func(event *stripe.Event) error{
		"payment_intent.requires_action": payment.HandlePaymentIntent(database.RequiresAction),
		"payment_intent.processing":      payment.HandlePaymentIntent(database.Processing),
		"payment_intent.succeeded":       payment.HandlePaymentSuccess,
		"payment_intent.payment_failed":  payment.HandlePaymentFailed,
		"payment_intent.canceled":        payment.HandlePaymentCancel,
		"refund.created":                 payment.HandleRefundCreated,
		"refund.updated":                 payment.HandleRefundUpdated,
		"refund.failed":                  payment.HandleRefundFailed,
		"charge.refunded":                payment.HandleChargeRefunded,
		"charge.dispute.created":         payment.HandleChargeDisputeCreated,
		"charge.dispute.closed":          payment.HandleChargeDisputeClosed,
	}

	if handler, has := handlers[string(event.Type)]; has {
		if err := handler(&event); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (payment *Payments) HandlePaymentIntent(status database.PaymentStatus) func(event *stripe.Event) error {
	return func(event *stripe.Event) error {
		var intent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
			payment.ErrorLog.Printf("Failed to unmarshal payment intent: %s\n", err)
			return errors.New("unexpected internal server error")
		}

		if paymentKey, exists := intent.Metadata["payment_key"]; exists {
			coursePurchase, err := payment.Database.GetCoursePurchaseByPaymentKey(paymentKey)
			if err != nil {
				payment.ErrorLog.Printf("Failed to get course purchase by payment key (\"%s\"): %s\n", paymentKey, err)
				return errors.New("unexpected internal server error")
			}

			payment.InfoLog.Println("Found course purchase by payment key")

			if err := payment.Database.UpdateCoursePurchasePaymentStatus(coursePurchase.ID, status); err != nil {
				payment.ErrorLog.Printf("Failed to update course purchase's (\"%s\") payment status: %s\n", coursePurchase.ID, err)
				return errors.New("unexpected internal server error")
			}

			payment.InfoLog.Println("Managed to update course payment status!")
		} else {
			payment.WarningLog.Println("Payment key wasn't found in meta data.")
		}

		return nil
	}
}

func (payment *Payments) HandlePaymentSuccess(event *stripe.Event) error {
	var intent stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
		payment.ErrorLog.Printf("Failed to unmarshal payment intent: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	if paymentKey, exists := intent.Metadata["payment_key"]; exists {
		coursePurchase, err := payment.Database.GetCoursePurchaseByPaymentKey(paymentKey)
		if err != nil {
			payment.ErrorLog.Printf("Failed to get course purchase by payment key (\"%s\"): %s\n", paymentKey, err)
			return errors.New("unexpected internal server error")
		}

		if err := payment.Database.UpdateCoursePurchasePaymentStatus(coursePurchase.ID, database.Succeeded); err != nil {
			payment.ErrorLog.Printf("Failed to update course purchase's (\"%s\") payment status: %s\n", coursePurchase.ID, err)
			return errors.New("unexpected internal server error")
		}

		if coursePurchase.AffiliateCode.Valid {
			affiliateUser, err := payment.Database.GetUserByAffiliateCode(coursePurchase.AffiliateCode.String, database.All)
			if err != nil {
				payment.ErrorLog.Printf("Failed to get user by affiliate code (\"%s\"): %s\n", coursePurchase.AffiliateCode.String, err)
				return errors.New("unexpected internal server error")
			}

			if err := payment.Database.RegisterAffiliatePointsChange(affiliateUser.ID, coursePurchase.CourseID, AffiliateReward, "Affiliate reward received"); err != nil {
				payment.ErrorLog.Printf("Failed to reward user (\"%s\") with affiliate points: %s\n", affiliateUser.ID, err)
				return errors.New("unexpected internal server error")
			}
		}

		// TODO: Send email to thank user for purchase.
	}

	return nil
}

func (payment *Payments) HandlePaymentCancel(event *stripe.Event) error {
	var intent stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
		payment.ErrorLog.Printf("Failed to unmarshal payment intent: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	if paymentKey, exists := intent.Metadata["payment_key"]; exists {
		coursePurchase, err := payment.Database.GetCoursePurchaseByPaymentKey(paymentKey)
		if err != nil {
			payment.ErrorLog.Printf("Failed to get course purchase by payment key (\"%s\"): %s\n", paymentKey, err)
			return errors.New("unexpected internal server error")
		}

		if err := payment.Database.UpdateCoursePurchasePaymentStatus(coursePurchase.ID, database.Cancelled); err != nil {
			payment.ErrorLog.Printf("Failed to update course purchase's (\"%s\") payment status: %s\n", coursePurchase.ID, err)
			return errors.New("unexpected internal server error")
		}

		if coursePurchase.AffiliatePointsUsed > 0 {
			if err := payment.Database.RegisterAffiliatePointsChange(coursePurchase.UserID, coursePurchase.CourseID, int(coursePurchase.AffiliatePointsUsed), "Payment cancelled"); err != nil {
				payment.ErrorLog.Printf("Failed to refund affiliate points after payment was cancelled: %s\n", err)
				return errors.New("unexpected internal server error")
			}
		}
	}

	return nil
}

func (payment *Payments) HandlePaymentFailed(event *stripe.Event) error {
	var intent stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
		payment.ErrorLog.Printf("Failed to unmarshal payment intent: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	if paymentKey, exists := intent.Metadata["payment_key"]; exists {
		coursePurchase, err := payment.Database.GetCoursePurchaseByPaymentKey(paymentKey)
		if err != nil {
			payment.ErrorLog.Printf("Failed to get course purchase by payment key (\"%s\"): %s\n", paymentKey, err)
			return errors.New("unexpected internal server error")
		}

		if err := payment.Database.UpdateCoursePurchasePaymentStatus(coursePurchase.ID, database.Failed); err != nil {
			payment.ErrorLog.Printf("Failed to update course purchase's (\"%s\") payment status: %s\n", coursePurchase.ID, err)
			return errors.New("unexpected internal server error")
		}

		if coursePurchase.AffiliatePointsUsed > 0 {
			if err := payment.Database.RegisterAffiliatePointsChange(coursePurchase.UserID, coursePurchase.CourseID, int(coursePurchase.AffiliatePointsUsed), "Payment failed"); err != nil {
				payment.ErrorLog.Printf("Failed to refund affiliate points after payment failed: %s\n", err)
				return errors.New("unexpected internal server error")
			}
		}
	}

	return nil
}

func (payment *Payments) HandleRefundCreated(event *stripe.Event) error {
	var refund stripe.Refund
	if err := json.Unmarshal(event.Data.Raw, &refund); err != nil {
		payment.ErrorLog.Printf("Failed to unmarshal refund: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	if len(refund.Metadata) == 0 {
		payment.ErrorLog.Println("Refund doesn't have any metadata")
		return errors.New("refund doesn't contain any metadata")
	}

	paymentKey, hasPaymentKey := refund.Metadata["payment_key"]
	coursePurchaseId, hasCoursePurchaseId := refund.Metadata["course_purchase_id"]

	if !(hasPaymentKey || hasCoursePurchaseId) {
		payment.ErrorLog.Println("Refund metadata didn't contain a payment key nor a course purchase id")
		return errors.New("refund doesn't contain required metadata")
	}

	var coursePurchase *models.CoursePurchaseModel
	var err error

	if hasPaymentKey {
		coursePurchase, err = payment.Database.GetCoursePurchaseByPaymentKey(paymentKey)
		if err != nil || coursePurchase == nil {
			payment.ErrorLog.Printf("Failed to find course purchase by payment key (\"%s\"): %s\n", paymentKey, err)
			return errors.New("unexpected internal server error")
		}
	} else if hasCoursePurchaseId {
		coursePurchase, err := payment.Database.GetCoursePurchaseByID(coursePurchaseId)
		if err != nil || coursePurchase == nil {
			payment.ErrorLog.Printf("Failed to find course purchase by ID (\"%s\"): %s\n", coursePurchaseId, err)
			return errors.New("unexpected internal server error")
		}
	}

	status := database.RefundPending
	switch refund.Status {
	case "pending":
		status = database.RefundPending
	case "requires_action":
		status = database.RefundRequiresAction
	case "succeeded":
		status = database.RefundSucceeded
	case "failed":
		status = database.RefundFailed
	case "canceled":
		status = database.RefundCancelled
	}

	if err := payment.Database.RegisterRefund(coursePurchase.UserID, coursePurchase.ID, status); err != nil {
		payment.ErrorLog.Printf("Failed to insert new refund: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	if err := payment.Database.UpdateCoursePurchasePaymentStatus(coursePurchase.ID, database.Refunded); err != nil {
		payment.ErrorLog.Printf("Failed to update course purchase (\"%s\") payment status to refunded: %s\n", coursePurchase.ID, err)
		return errors.New("unexpected internal server error")
	}

	return nil
}

func (payment *Payments) HandleRefundUpdated(event *stripe.Event) error {
	var refund stripe.Refund
	if err := json.Unmarshal(event.Data.Raw, &refund); err != nil {
		payment.ErrorLog.Printf("Failed to unmarshal refund: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	payment.InfoLog.Printf("Refund updated with status: %s\n", refund.Status)
	payment.InfoLog.Println("Refund updated with the following metadata:")

	for key, value := range refund.Metadata {
		payment.InfoLog.Printf("%s : %s\n", key, value)
	}

	// TODO: Update refund request to refund.Status.
	return nil
}

func (payment *Payments) HandleRefundFailed(event *stripe.Event) error {
	var refund stripe.Refund
	if err := json.Unmarshal(event.Data.Raw, &refund); err != nil {
		payment.ErrorLog.Printf("Failed to unmarshal refund: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	payment.InfoLog.Printf("Refund failed with status: %s\n", refund.Status)
	payment.InfoLog.Println("Refund failed with the following metadata:")

	for key, value := range refund.Metadata {
		payment.InfoLog.Printf("%s : %s\n", key, value)
	}

	// TODO: Update refund request to "Failed".
	return nil
}

func (payment *Payments) HandleChargeRefunded(event *stripe.Event) error {
	var charge stripe.Charge
	if err := json.Unmarshal(event.Data.Raw, &charge); err != nil {
		payment.ErrorLog.Printf("Failed to unmarshal charge: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	payment.InfoLog.Println("Charge refunded with refund status: ", charge.Refunded)
	payment.InfoLog.Println("Charge refunds: ", charge.Refunds)

	if len(charge.Metadata) > 0 {
		payment.InfoLog.Println("Charge refund with the following metadata:")
		for key, value := range charge.Metadata {
			payment.InfoLog.Printf("%s : %s\n", key, value)
		}
	} else {
		payment.InfoLog.Println("Charge refund doesn't have any metadata:")
	}

	if charge.PaymentIntent != nil {
		if len(charge.PaymentIntent.Metadata) > 0 {
			payment.InfoLog.Println("Charge refund's payment intent metadata:")
			for key, value := range charge.PaymentIntent.Metadata {
				payment.InfoLog.Printf("%s : %s\n", key, value)
			}
		} else {
			payment.InfoLog.Println("Charge refund payment intent doesn't have any metadata")
		}
	} else {
		payment.InfoLog.Println("Charge refund doesn't have a payment intent")
	}

	// TODO: Update refund request to "Refunded".
	return nil
}

func (payment *Payments) HandleChargeDisputeCreated(event *stripe.Event) error {
	var dispute stripe.Dispute
	if err := json.Unmarshal(event.Data.Raw, &dispute); err != nil {
		payment.ErrorLog.Printf("Failed to unmarshal dispute: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	return nil
}

func (payment *Payments) HandleChargeDisputeClosed(event *stripe.Event) error {
	var dispute stripe.Dispute
	if err := json.Unmarshal(event.Data.Raw, &dispute); err != nil {
		payment.ErrorLog.Printf("Failed to unmarshal dispute: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	return nil
}
