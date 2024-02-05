package domain

import "time"

type User struct {
	ID       int64
	Email    string
	Password string
	Ctime    time.Time
}

type UserProfile struct {
	ID       int64
	Nickname string
	Birthday time.Time
	AboutMe  string
}
