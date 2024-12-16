package models

import "time"

type CertificateModel struct {
	ID        string
	UserID    string
	CourseID  string
	CreatedAt time.Time
}
