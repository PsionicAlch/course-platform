package payments

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
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

	if paymentKey, hasPaymentKey := refund.Metadata["payment_key"]; hasPaymentKey {
		coursePurchase, err := payment.Database.GetCoursePurchaseByPaymentKey(paymentKey)
		if err != nil || coursePurchase == nil {
			payment.ErrorLog.Printf("Failed to find course purchase by payment key (\"%s\"): %s\n", paymentKey, err)
			return errors.New("unexpected internal server error")
		}

		if err := payment.Database.RegisterRefund(coursePurchase.UserID, coursePurchase.ID, database.RefundPending); err != nil {
			payment.ErrorLog.Printf("Failed to insert new refund: %s\n", err)
			return errors.New("unexpected internal server error")
		}

		if err := payment.Database.UpdateCoursePurchasePaymentStatus(coursePurchase.ID, database.Refunded); err != nil {
			payment.ErrorLog.Printf("Failed to update course purchase (\"%s\") payment status to refunded: %s\n", coursePurchase.ID, err)
			return errors.New("unexpected internal server error")
		}

		return nil
	} else {
		payment.ErrorLog.Println("Refund metadata didn't contain a payment key")
		return errors.New("refund doesn't contain required metadata")
	}
}

func (payment *Payments) HandleRefundUpdated(event *stripe.Event) error {
	var refund stripe.Refund
	if err := json.Unmarshal(event.Data.Raw, &refund); err != nil {
		payment.ErrorLog.Printf("Failed to unmarshal refund: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	if paymentKey, hasPaymentKey := refund.Metadata["payment_key"]; hasPaymentKey {
		coursePurchase, err := payment.Database.GetCoursePurchaseByPaymentKey(paymentKey)
		if err != nil || coursePurchase == nil {
			payment.ErrorLog.Printf("Failed to find course purchase by payment key (\"%s\"): %s\n", paymentKey, err)
			return errors.New("unexpected internal server error")
		}

		status := database.RefundPending
		switch refund.Status {
		case "requires_action":
			status = database.RefundRequiresAction
		case "canceled":
			status = database.RefundCancelled
		default:
			return nil
		}

		if status == database.RefundFailed || status == database.RefundCancelled {
			if err := payment.Database.UpdateCoursePurchasePaymentStatus(coursePurchase.ID, database.Succeeded); err != nil {
				payment.ErrorLog.Printf("Failed to update course purchase (\"%s\") payment status to succeeded: %s\n", coursePurchase.ID, err)
				return errors.New("unexpected internal server error")
			}
		}

		refundModel, err := payment.Database.GetRefundWithCoursePurchaseID(coursePurchase.ID)
		if err != nil {
			payment.ErrorLog.Printf("Failed to get refund from course purchase ID (\"%s\"): %s\n", coursePurchase.ID, err)
			return errors.New("unexpected internal server error")
		}

		if err := payment.Database.UpdateRefundStatus(coursePurchase.ID, status); err != nil {
			payment.ErrorLog.Printf("Failed to update refund (\"%s\") status: %s\n", refundModel.ID, err)
			return errors.New("unexpected internal server error")
		}

		return nil
	} else {
		payment.ErrorLog.Println("Refund metadata didn't contain a payment key")
		return errors.New("refund doesn't contain required metadata")
	}
}

func (payment *Payments) HandleRefundFailed(event *stripe.Event) error {
	var refund stripe.Refund
	if err := json.Unmarshal(event.Data.Raw, &refund); err != nil {
		payment.ErrorLog.Printf("Failed to unmarshal refund: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	if paymentKey, hasPaymentKey := refund.Metadata["payment_key"]; hasPaymentKey {
		coursePurchase, err := payment.Database.GetCoursePurchaseByPaymentKey(paymentKey)
		if err != nil || coursePurchase == nil {
			payment.ErrorLog.Printf("Failed to find course purchase by payment key (\"%s\"): %s\n", paymentKey, err)
			return errors.New("unexpected internal server error")
		}

		refundModel, err := payment.Database.GetRefundWithCoursePurchaseID(coursePurchase.ID)
		if err != nil {
			payment.ErrorLog.Printf("Failed to get refund from course purchase ID (\"%s\"): %s\n", coursePurchase.ID, err)
			return errors.New("unexpected internal server error")
		}

		if err := payment.Database.UpdateRefundStatus(coursePurchase.ID, database.RefundFailed); err != nil {
			payment.ErrorLog.Printf("Failed to update refund (\"%s\") status to failed: %s\n", refundModel.ID, err)
			return errors.New("unexpected internal server error")
		}

		if err := payment.Database.UpdateCoursePurchasePaymentStatus(coursePurchase.ID, database.Succeeded); err != nil {
			payment.ErrorLog.Printf("Failed to update course purchase (\"%s\") payment status to succeeded: %s\n", coursePurchase.ID, err)
			return errors.New("unexpected internal server error")
		}

		return nil
	} else {
		payment.ErrorLog.Println("Refund metadata didn't contain a payment key")
		return errors.New("refund doesn't contain required metadata")
	}
}

func (payment *Payments) HandleChargeRefunded(event *stripe.Event) error {
	var charge stripe.Charge
	if err := json.Unmarshal(event.Data.Raw, &charge); err != nil {
		payment.ErrorLog.Printf("Failed to unmarshal charge: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	if paymentKey, hasPaymentKey := charge.Metadata["payment_key"]; hasPaymentKey {
		coursePurchase, err := payment.Database.GetCoursePurchaseByPaymentKey(paymentKey)
		if err != nil || coursePurchase == nil {
			payment.ErrorLog.Printf("Failed to find course purchase by payment key (\"%s\"): %s\n", paymentKey, err)
			return errors.New("unexpected internal server error")
		}

		refundModel, err := payment.Database.GetRefundWithCoursePurchaseID(coursePurchase.ID)
		if err != nil {
			payment.ErrorLog.Printf("Failed to get refund from course purchase ID (\"%s\"): %s\n", coursePurchase.ID, err)
			return errors.New("unexpected internal server error")
		}

		if err := payment.Database.UpdateRefundStatus(refundModel.ID, database.RefundSucceeded); err != nil {
			payment.ErrorLog.Printf("Failed to update refund (\"%s\") status to succeeded: %s\n", refundModel.ID, err)
			return errors.New("unexpected internal server error")
		}

		return nil
	} else {
		payment.ErrorLog.Println("Charge metadata didn't contain a payment key")
		return errors.New("charge doesn't contain required metadata")
	}
}

func (payment *Payments) HandleChargeDisputeCreated(event *stripe.Event) error {
	var dispute stripe.Dispute
	if err := json.Unmarshal(event.Data.Raw, &dispute); err != nil {
		payment.ErrorLog.Printf("Failed to unmarshal dispute: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	if paymentKey, hasPaymentKey := dispute.Metadata["payment_key"]; hasPaymentKey {
		coursePurchase, err := payment.Database.GetCoursePurchaseByPaymentKey(paymentKey)
		if err != nil || coursePurchase == nil {
			payment.ErrorLog.Printf("Failed to find course purchase by payment key (\"%s\"): %s\n", paymentKey, err)
			return errors.New("unexpected internal server error")
		}

		var status database.RefundStatus
		switch dispute.Status {
		case "warning_needs_response":
			status = database.DisputeWarningNeedsResponse
		case "warning_under_review":
			status = database.DisputeWarningUnderReview
		case "warning_closed":
			status = database.DisputeWarningClosed
		case "needs_response":
			status = database.DisputeNeedsResponse
		case "under_review":
			status = database.DisputeUnderReview
		case "won":
			status = database.DisputeWon
		case "lost":
			status = database.DisputeLost
		}

		if err := payment.Database.RegisterRefund(coursePurchase.UserID, coursePurchase.ID, status); err != nil {
			payment.ErrorLog.Printf("Failed to insert new dispute: %s\n", err)
			return errors.New("unexpected internal server error")
		}

		if !slices.Contains([]database.RefundStatus{database.DisputeWarningClosed, database.DisputeWon}, status) {
			if err := payment.Database.UpdateCoursePurchasePaymentStatus(coursePurchase.ID, database.Disputed); err != nil {
				payment.ErrorLog.Printf("Failed to update course purchase (\"%s\") payment status to disputed: %s\n", coursePurchase.ID, err)
				return errors.New("unexpected internal server error")
			}
		}

		return nil
	} else {
		payment.ErrorLog.Println("Dispute metadata didn't contain a payment key")
		return errors.New("charge doesn't contain required metadata")
	}
}

func (payment *Payments) HandleChargeDisputeClosed(event *stripe.Event) error {
	var dispute stripe.Dispute
	if err := json.Unmarshal(event.Data.Raw, &dispute); err != nil {
		payment.ErrorLog.Printf("Failed to unmarshal dispute: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	if paymentKey, hasPaymentKey := dispute.Metadata["payment_key"]; hasPaymentKey {
		coursePurchase, err := payment.Database.GetCoursePurchaseByPaymentKey(paymentKey)
		if err != nil || coursePurchase == nil {
			payment.ErrorLog.Printf("Failed to find course purchase by payment key (\"%s\"): %s\n", paymentKey, err)
			return errors.New("unexpected internal server error")
		}

		refundModel, err := payment.Database.GetRefundWithCoursePurchaseID(coursePurchase.ID)
		if err != nil {
			payment.ErrorLog.Printf("Failed to get refund from course purchase ID (\"%s\"): %s\n", coursePurchase.ID, err)
			return errors.New("unexpected internal server error")
		}

		var status database.RefundStatus
		switch dispute.Status {
		case "warning_needs_response":
			status = database.DisputeWarningNeedsResponse
		case "warning_under_review":
			status = database.DisputeWarningUnderReview
		case "warning_closed":
			status = database.DisputeWarningClosed
		case "needs_response":
			status = database.DisputeNeedsResponse
		case "under_review":
			status = database.DisputeUnderReview
		case "won":
			status = database.DisputeWon
		case "lost":
			status = database.DisputeLost
		}

		if err := payment.Database.UpdateRefundStatus(refundModel.ID, status); err != nil {
			payment.ErrorLog.Printf("Failed to update dispute status: %s\n", err)
			return errors.New("unexpected internal server error")
		}

		if slices.Contains([]database.RefundStatus{database.DisputeWarningClosed, database.DisputeWon}, status) {
			if err := payment.Database.UpdateCoursePurchasePaymentStatus(coursePurchase.ID, database.Succeeded); err != nil {
				payment.ErrorLog.Printf("Failed to update course purchase (\"%s\") payment status to succeeded: %s\n", coursePurchase.ID, err)
				return errors.New("unexpected internal server error")
			}
		}

		return nil
	} else {
		payment.ErrorLog.Println("Dispute metadata didn't contain a payment key")
		return errors.New("charge doesn't contain required metadata")
	}
}
