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

	payment.InfoLog.Printf("Processing event: %s\n", event.Type)

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

		if err := payment.HandlePaymentIntent(&event, database.Succeeded); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// TODO: Send email to thank the user for the successful purchase.
	case "payment_intent.payment_failed":
		payment.InfoLog.Println("Handling event: ", "payment_intent.payment_failed")

		if err := payment.HandlePaymentIntent(&event, database.Failed); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case "payment_intent.canceled":
		payment.InfoLog.Println("Handling event: ", "payment_intent.canceled")

		if err := payment.HandlePaymentIntent(&event, database.Cancelled); err != nil {
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
	default:
		payment.WarningLog.Printf("Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

func (payment *Payments) HandlePaymentIntent(event *stripe.Event, status database.PaymentStatus) error {
	var intent stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
		payment.ErrorLog.Printf("Failed to unmarshal payment intent: %s\n", err)
		return errors.New("unexpected internal server error")
	}

	payment.InfoLog.Println("Handling payment intent")

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
