package emails

import (
	"bytes"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/email"
	smtp_email_client "github.com/PsionicAlch/psionicalch-home/internal/email/clients/SMTP"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/website/config"
	"github.com/PsionicAlch/psionicalch-home/website/html"
)

type Emails struct {
	utils.Loggers
	Client email.EmailClient
	Render render.Renderer
}

func SetupEmails(emailRenderer render.Renderer) *Emails {
	loggers := utils.CreateLoggers("EMAIL")

	emailProvider := config.GetWithoutError[string]("EMAIL_PROVIDER")
	emailHost := config.GetWithoutError[string]("EMAIL_HOST")
	emailPort := config.GetWithoutError[string]("EMAIL_PORT")
	emailAddr := config.GetWithoutError[string]("EMAIL_ADDRESS")
	emailPassword := config.GetWithoutError[string]("EMAIL_PASSWORD")

	var client email.EmailClient

	switch emailProvider {
	case "smtp":
		client = smtp_email_client.SetupSMTPEmailClient(emailHost, emailPort, emailAddr, emailPassword)
	default:
		loggers.ErrorLog.Fatalf("%s is not a valid email provider", emailProvider)
	}

	return &Emails{
		Loggers: loggers,
		Client:  client,
		Render:  emailRenderer,
	}
}

func (e *Emails) SendWelcomeEmail(email, firstName, affiliateCode string) {
	emailData := html.NewGreetingEmail(firstName, affiliateCode)

	buf := new(bytes.Buffer)
	if err := e.Render.Render(buf, "greeting", emailData); err != nil {
		e.ErrorLog.Printf("Failed to render greeting email for %s: %s\n", email, err)
		return
	}

	e.Client.SendEmail(email, emailData.Title, buf.String())
}

func (e *Emails) SendLoginEmail(email, firstName, ipAddr string, date time.Time) {
	emailData := html.NewLoginEmail(firstName, ipAddr, date)

	buf := new(bytes.Buffer)
	if err := e.Render.Render(buf, "login", emailData); err != nil {
		e.ErrorLog.Printf("Failed to render login email for %s: %s\n", email, err)
		return
	}

	e.Client.SendEmail(email, emailData.Title, buf.String())
}

func (e *Emails) SendPasswordResetEmail(email, firstName, emailToken string) {
	emailData := html.NewPasswordResetEmail(firstName, emailToken)

	buf := new(bytes.Buffer)
	if err := e.Render.Render(buf, "reset-password", emailData); err != nil {
		e.ErrorLog.Printf("Failed to render reset-password email for %s: %s\n", email, err)
		return
	}

	e.Client.SendEmail(email, emailData.Title, buf.String())
}
