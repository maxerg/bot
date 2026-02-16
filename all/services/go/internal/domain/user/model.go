package user

import "time"

type User struct {
	ID        int64
	TgUserID  int64
	Username  string
	FirstName string
	LastName  string
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
