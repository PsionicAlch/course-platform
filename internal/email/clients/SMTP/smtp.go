package smtp_email_client

import (
	"net/smtp"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
)

type SMTPEmailClient struct {
	utils.Loggers
	Host        string
	Port        string
	SenderEmail string
	Password    string
}

func SetupSMTPEmailClient(host, port, email, password string) *SMTPEmailClient {
	loggers := utils.CreateLoggers("SMTP EMAIL CLIENT")

	return &SMTPEmailClient{
		Loggers:     loggers,
		Host:        host,
		Port:        port,
		SenderEmail: email,
		Password:    password,
	}
}

func (client *SMTPEmailClient) SendEmail(recipient, subject, body string) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte("Subject: " + subject + "\n" + mime + body)

	auth := smtp.PlainAuth("", client.SenderEmail, client.Password, client.Host)
	err := smtp.SendMail(client.Host+":"+client.Port, auth, client.SenderEmail, []string{recipient}, msg)
	if err != nil {
		client.ErrorLog.Printf("Failed to send email to %s: %s\n", recipient, err)
	}
}
