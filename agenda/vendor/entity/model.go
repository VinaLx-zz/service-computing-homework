package entity

type Model interface {
	CancelMeeting(user *User, title string) bool
}
