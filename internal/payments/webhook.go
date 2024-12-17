package payments

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

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

	switch event.Type {
	case "payment_intent.requires_action":
		payment.InfoLog.Println("Handling event: ", "payment_intent.requires_action")

		if err := payment.HandlePaymentIntent(&event, database.RequiresAction); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case "payment_intent.processing":
		payment.InfoLog.Println("Handling event: ", "payment_intent.processing")

		if err := payment.HandlePaymentIntent(&event, database.Processing); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case "payment_intent.succeeded":
		payment.InfoLog.Println("Handling event: ", "payment_intent.succeeded")

		if err := payment.HandlePaymentSuccess(&event); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case "payment_intent.payment_failed":
		payment.InfoLog.Println("Handling event: ", "payment_intent.payment_failed")

		if err := payment.HandlePaymentFailed(&event); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case "payment_intent.canceled":
		payment.InfoLog.Println("Handling event: ", "payment_intent.canceled")

		if err := payment.HandlePaymentCancel(&event); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case "refund.created":
		var refund stripe.Refund
		if err := json.Unmarshal(event.Data.Raw, &refund); err != nil {
			payment.ErrorLog.Printf("Failed to unmarshal refund: %s\n", err)
			http.Error(w, "Failed to parse webhook", http.StatusBadRequest)
			return
		}

		// TODO: Create or update refund request.
	case "refund.updated":
		var refund stripe.Refund
		if err := json.Unmarshal(event.Data.Raw, &refund); err != nil {
			payment.ErrorLog.Printf("Failed to unmarshal refund: %s\n", err)
			http.Error(w, "Failed to parse webhook", http.StatusBadRequest)
			return
		}

		// TODO: Update refund request to refund.Status.
	case "refund.failed":
		var refund stripe.Refund
		if err := json.Unmarshal(event.Data.Raw, &refund); err != nil {
			payment.ErrorLog.Printf("Failed to unmarshal refund: %s\n", err)
			http.Error(w, "Failed to parse webhook", http.StatusBadRequest)
			return
		}

		// TODO: Update refund request to "Failed".
	}

	w.WriteHeader(http.StatusOK)
}

func (payment *Payments) HandlePaymentIntent(event *stripe.Event, status database.PaymentStatus) error {
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
