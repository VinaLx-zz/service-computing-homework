package user

import "time"

var count uint64

// User type
type User struct {
	UID        uint64
	Username   string
	Password   string
	SignUpTime time.Time
}

// ResetUID count
func ResetUID(i uint64) {
	count = i
}

// NewUser returns a new user with a new uid
func NewUser(username, password string) *User {
	u := User{
		UID:        count,
		Username:   username,
		Password:   password,
		SignUpTime: time.Now(),
	}
	count++
	return &u
}
