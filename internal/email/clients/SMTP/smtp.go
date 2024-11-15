package smtp_email_client

type SMTPEmailClient struct {
	Host        string
	Port        string
	SenderEmail string
	Password    string
}

func (client *SMTPEmailClient) SendTextEmail(recipient string) {
}
