package html

import (
	"net/url"
	"time"
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
	LatestCourses []struct {
		Name        string
		Description string
		Slug        string
	}
	AffiliateCode string
}

func NewGreetingEmail(firstName, affiliateCode string) *GreetingEmail {
	return &GreetingEmail{
		BaseEmail: NewBaseEmail("Welcome to PsionicAlch"),
		FirstName: firstName,
		LatestCourses: []struct {
			Name        string
			Description string
			Slug        string
		}{},
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
