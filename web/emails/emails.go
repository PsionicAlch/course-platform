package emails

import (
	"bytes"
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/internal/email"
	smtp_email_client "github.com/PsionicAlch/psionicalch-home/internal/email/clients/smtp"
	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/web/config"
	"github.com/PsionicAlch/psionicalch-home/web/html"
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

func (e *Emails) SendEmail(email, title, tmpl string, data any) {
	buf := new(bytes.Buffer)
	if err := e.Render.Render(buf, nil, tmpl, data); err != nil {
		e.ErrorLog.Printf("Failed to \"%s\" email for %s: %s\n", tmpl, email, err)
		return
	}

	e.Client.SendEmail(email, title, buf.String())
}

func (e *Emails) SendWelcomeEmail(email, firstName, affiliateCode string, discount *models.DiscountModel, latestCourses []*models.CourseModel) {
	emailData := html.NewGreetingEmail(firstName, affiliateCode, discount, latestCourses)
	e.SendEmail(email, emailData.Title, "greeting", emailData)
}

func (e *Emails) SendLoginEmail(email, firstName, ipAddr string, date time.Time) {
	emailData := html.NewLoginEmail(firstName, ipAddr, date)
	e.SendEmail(email, emailData.Title, "login", emailData)
}

func (e *Emails) SendPasswordResetEmail(email, firstName, emailToken string) {
	emailData := html.NewPasswordResetEmail(firstName, emailToken)
	e.SendEmail(email, emailData.Title, "reset-password", emailData)
}

func (e *Emails) SendPasswordResetConfirmationEmail(email, firstName string) {
	emailData := html.NewPasswordResetConfirmationEmail(firstName)
	e.SendEmail(email, emailData.Title, "password-reset-confirmation", emailData)
}

func (e *Emails) SendSuspiciousActivityEmail(email, firstName, ipAddr string, dateTime time.Time) {
	emailData := html.NewSuspiciousActivityEmail(firstName, ipAddr, dateTime)
	e.SendEmail(email, emailData.Title, "suspicious-activity", emailData)
}

func (e *Emails) SendAccountDeletionEmail(email, firstName string) {
	emailData := html.NewAccountDeletionEmail(firstName)
	e.SendEmail(email, emailData.Title, "account-deletion", emailData)
}

func (e *Emails) SendRefundRequestAcknowledgementEmail(email, firstName string) {
	emailData := html.NewRefundRequestAcknowledgementEmail(firstName)
	e.SendEmail(email, emailData.Title, "refund-request-acknowledgement", emailData)
}

func (e *Emails) SendThankYouForPurchaseEmail(email, firstName, affiliateCode string, course *models.CourseModel, discount *models.DiscountModel) {
	emailData := html.NewThankYouForPurchaseEmail(firstName, affiliateCode, course, discount)
	e.SendEmail(email, emailData.Title, "thank-you-for-purchase", emailData)
}

func (e *Emails) SendRefundRequestFailedEmail(email, firstName, courseName, failureReason string) {
	emailData := html.NewRefundRequestFailedEmail(firstName, courseName, failureReason)
	e.SendEmail(email, emailData.Title, "refund-request-failed", emailData)
}

func (e *Emails) SendRefundRequestCancelledEmail(email, firstName, courseName string) {
	emailData := html.NewRefundRequestCancelledEmail(firstName, courseName)
	e.SendEmail(email, emailData.Title, "refund-request-cancelled", emailData)
}

func (e *Emails) SendRefundRequestSuccessfulEmail(email, firstName, courseName string, refundAmount float64) {
	emailData := html.NewRefundRequestSuccessfulEmail(firstName, courseName, refundAmount)
	e.SendEmail(email, emailData.Title, "refund-request-successful", emailData)
}
