package model

import (
	"crypto/md5"
	"encoding/json"
	"entity"
	"err"
	"os"
)

func openFileRewrite(path string) (*os.File, error) {
	return os.OpenFile(path, rewritePerm, 0644)
}

const rewritePerm = os.O_WRONLY | os.O_CREATE | os.O_TRUNC

// LoadUsers from UserFile()
func LoadUsers() entity.Users {
	file, e := os.Open(UserFile())
	if e != nil {
		return entity.NewUsers()
	}
	us, e := entity.DeserializeUser(file)
	err.LogFatalIfError(e)
	return us
}

// LoadMeetings from Meeting
func LoadMeetings() *entity.Meetings {
	file, err := os.Open(MeetingFile())
	if err != nil {
		return entity.NewMeetings()
	}
	ms, err := entity.DeserializeMeeting(file)
	return ms
}

type login struct {
	Username, Password string
}

func loadLoginFile() *login {
	file, e := os.Open(LoginFile())
	if e != nil {
		return nil
	}
	l := new(login)
	e = json.NewDecoder(file).Decode(l)
	if e != nil {
		Logout()
		return nil
	}
	return l
}

func writeLoginFile(user, pass string) {
	file, e := openFileRewrite(LoginFile())
	err.LogFatalIfError(e)
	json.NewEncoder(file).Encode(login{user, pass})
}

func validPassword(u *entity.User, pass string) bool {
	return string(md5.New().Sum([]byte(pass))) == u.Password
}

// LoadLogin file to get current login user
// returns nil if passwords don't match the hash in users or
// login file doesn't exist
func LoadLogin(users entity.Users) (*entity.User, err.Err) {
	l := loadLoginFile()
	if l == nil {
		return nil, err.OK
	}
	u := users.Lookup(l.Username)
	if u == nil {
		return nil, err.NoSuchUser
	}
	if !validPassword(u, l.Password) {
		return nil, err.InconsistentState
	}
	return u, err.OK
}

// Logout try to delete the login file returns true if success
// returns false if there's no login file
func Logout() bool {
	return os.Remove(LoginFile()) == nil
}

// Login takes user set and the user and password to login
// and writes the login state to file if success
func Login(users entity.Users, user, pass string) err.Err {
	u := users.Lookup(user)
	if u == nil {
		return err.NoSuchUser
	}
	if validPassword(u, pass) {
		writeLoginFile(user, pass)
		return err.OK
	}
	return err.AuthenticateFail
}

// StoreUser to UserFile
func StoreUser(users entity.Users) {
	file, e := openFileRewrite(UserFile())
	err.LogFatalIfError(e)
	users.Serialize(file)
}

// StoreMeeting to MeetingFile
func StoreMeeting(meetings entity.Meetings) {
	file, e := openFileRewrite(MeetingFile())
	err.LogFatalIfError(e)
	meetings.Serialize(file)
}
