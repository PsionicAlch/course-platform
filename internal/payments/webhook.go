package payments

import (
	"io"
	"net/http"

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
	default:
		payment.WarningLog.Printf("Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}
