package user

import (
	"time"
)

var count uint64

// Dao is User Data Access Object
type Dao interface {
	StoreUser(u *User) error
	GetAllUsers() ([]*User, error)
	GetUser(uid uint64) (*User, error)
}

// User type
type User struct {
	ID         uint64
	Username   string    `gorm:"not null"`
	Password   string    `gorm:"not null"`
	SignUpDate time.Time `gorm:"not null"`
}

// ResetUID count
func ResetUID(i uint64) {
	count = i
}

// NewUser returns a new user with a new uid
func NewUser(username, password string) *User {
	u := User{
		ID:         count,
		Username:   username,
		Password:   password,
		SignUpDate: time.Now(),
	}
	count++
	return &u
}
