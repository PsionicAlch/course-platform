package payments

import (
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

const PaymentToken = "payment"

// ValidatePaymentToken ensures that the provided payment token is valid.
func (payment *Payments) ValidatePaymentToken(token string) bool {
	paymentToken, err := payment.Database.GetToken(token, PaymentToken)
	if err != nil {
		payment.ErrorLog.Printf("Failed to get payment token from the database: %s\n", err)
		return false
	}

	if paymentToken == nil {
		return false
	}

	if paymentToken.Token == "" || paymentToken.TokenType != PaymentToken {
		return false
	}

	if time.Now().After(paymentToken.ValidUntil) {
		return false
	}

	return true
}

// GetUserFromPaymentToken gets the user from the provided payment token.
func (payment *Payments) GetUserFromPaymentToken(token string) (*models.UserModel, error) {
	return payment.Database.GetUserByToken(token, PaymentToken, database.All)
}

// DeletePaymentToken deletes the provided payment token from the database.
func (payment *Payments) DeletePaymentToken(token string) error {
	return payment.Database.DeleteToken(token, PaymentToken)
}
