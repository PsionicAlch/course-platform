package email

type EmailClient interface {
	SendEmail(recipient, subject, body string)
}
