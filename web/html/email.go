package html

import (
	"net/url"
	"time"

	"github.com/PsionicAlch/course-platform/internal/database/models"
)

type BaseEmail struct {
	Title string
}

func NewBaseEmail(title string) BaseEmail {
	return BaseEmail{
		Title: title,
	}
}

type GreetingEmail struct {
	BaseEmail
	FirstName     string
	Discount      *models.DiscountModel
	LatestCourses []*models.CourseModel
	AffiliateCode string
}

func NewGreetingEmail(firstName, affiliateCode string, discount *models.DiscountModel, latestCourses []*models.CourseModel) *GreetingEmail {
	return &GreetingEmail{
		BaseEmail:     NewBaseEmail("Welcome to PsionicAlch"),
		FirstName:     firstName,
		Discount:      discount,
		LatestCourses: latestCourses,
		AffiliateCode: affiliateCode,
	}
}

type LoginEmail struct {
	BaseEmail
	FirstName        string
	IPAddress        string
	URLSafeIPAddress string
	LoginDateTime    time.Time
}

func NewLoginEmail(firstName, ipAddr string, loginDateTime time.Time) *LoginEmail {
	return &LoginEmail{
		BaseEmail:        NewBaseEmail("Account Login Detected"),
		FirstName:        firstName,
		IPAddress:        ipAddr,
		URLSafeIPAddress: url.QueryEscape(ipAddr),
		LoginDateTime:    loginDateTime,
	}
}

type PasswordResetEmail struct {
	BaseEmail
	FirstName  string
	EmailToken string
}

func NewPasswordResetEmail(firstName, emailToken string) *PasswordResetEmail {
	return &PasswordResetEmail{
		BaseEmail:  NewBaseEmail("Password Reset Instructions"),
		FirstName:  firstName,
		EmailToken: emailToken,
	}
}

type PasswordResetConfirmationEmail struct {
	BaseEmail
	FirstName string
}

func NewPasswordResetConfirmationEmail(firstName string) *PasswordResetConfirmationEmail {
	return &PasswordResetConfirmationEmail{
		BaseEmail: NewBaseEmail("Password Update Confirmation"),
		FirstName: firstName,
	}
}

type SuspiciousActivityEmail struct {
	BaseEmail
	FirstName        string
	IPAddress        string
	URLSafeIPAddress string
	LoginDateTime    time.Time
}

func NewSuspiciousActivityEmail(firstName, ipAddr string, dateTime time.Time) *SuspiciousActivityEmail {
	return &SuspiciousActivityEmail{
		BaseEmail:        NewBaseEmail("Suspicious Account Activity"),
		FirstName:        firstName,
		IPAddress:        ipAddr,
		URLSafeIPAddress: url.QueryEscape(ipAddr),
		LoginDateTime:    dateTime,
	}
}

type AccountDeletionEmail struct {
	BaseEmail
	FirstName string
}

func NewAccountDeletionEmail(firstName string) *AccountDeletionEmail {
	return &AccountDeletionEmail{
		BaseEmail: NewBaseEmail("Account Deleted"),
		FirstName: firstName,
	}
}

type RefundRequestAcknowledgementEmail struct {
	BaseEmail
	FirstName string
}

func NewRefundRequestAcknowledgementEmail(firstName string) *RefundRequestAcknowledgementEmail {
	return &RefundRequestAcknowledgementEmail{
		BaseEmail: NewBaseEmail("Acknowledgement of Refund Request"),
		FirstName: firstName,
	}
}

type ThankYouForPurchaseEmail struct {
	BaseEmail
	FirstName     string
	AffiliateCode string
	Course        *models.CourseModel
	Discount      *models.DiscountModel
}

func NewThankYouForPurchaseEmail(firstName, affiliateCode string, course *models.CourseModel, discount *models.DiscountModel) *ThankYouForPurchaseEmail {
	return &ThankYouForPurchaseEmail{
		BaseEmail:     NewBaseEmail("Thank You For Your Purchase"),
		FirstName:     firstName,
		AffiliateCode: affiliateCode,
		Course:        course,
		Discount:      discount,
	}
}

type RefundRequestFailedEmail struct {
	BaseEmail
	FirstName     string
	CourseName    string
	FailureReason string
}

func NewRefundRequestFailedEmail(firstName, courseName, failureReason string) *RefundRequestFailedEmail {
	return &RefundRequestFailedEmail{
		BaseEmail:     NewBaseEmail("Refund Request Failed"),
		FirstName:     firstName,
		CourseName:    courseName,
		FailureReason: failureReason,
	}
}

type RefundRequestCancelledEmail struct {
	BaseEmail
	FirstName  string
	CourseName string
}

func NewRefundRequestCancelledEmail(firstName, courseName string) *RefundRequestCancelledEmail {
	return &RefundRequestCancelledEmail{
		BaseEmail:  NewBaseEmail("Refund Request Cancelled"),
		FirstName:  firstName,
		CourseName: courseName,
	}
}

type RefundRequestSuccessfulEmail struct {
	BaseEmail
	FirstName    string
	CourseName   string
	RefundAmount float64
}

func NewRefundRequestSuccessfulEmail(firstName, courseName string, refundAmount float64) *RefundRequestSuccessfulEmail {
	return &RefundRequestSuccessfulEmail{
		BaseEmail:    NewBaseEmail("Refund Request Successful"),
		FirstName:    firstName,
		CourseName:   courseName,
		RefundAmount: refundAmount,
	}
}
