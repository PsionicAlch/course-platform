package models

import "time"

type UserModel struct {
	ID         string
	Email      string
	Created_At time.Time
	Update_At  time.Time
}
