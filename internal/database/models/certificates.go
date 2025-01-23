package models

import "time"

// CertificateModel is a struct representation of the certificates table.
type CertificateModel struct {
	ID        string
	UserID    string
	CourseID  string
	CreatedAt time.Time
}
