package email

// Functions provided by the email package.
type EmailClient interface {
	SendEmail(recipient, subject, body string)
}
