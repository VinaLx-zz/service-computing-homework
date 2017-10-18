package model

import (
	"crypto/md5"
	"encoding/json"
	"entity"
	"err"
	"log"
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
		log.Printf("No User file; %s\n", e.Error())
		return entity.NewUsers()
	}
	us, e := entity.DeserializeUser(file)
	err.LogFatalIfError(e)
	return us
}

// LoadMeetings from Meeting
func LoadMeetings() *entity.Meetings {
	file, e := os.Open(MeetingFile())
	if e != nil {
		log.Printf("No Meeting file; %s\n", e.Error())
		return entity.NewMeetings()
	}
	ms, e := entity.DeserializeMeeting(file)
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
		log.Println("Login file not found, not logged in")
		return nil, err.OK
	}
	u := users.Lookup(l.Username)
	if u == nil {
		log.Printf("Login user not found: %s\n", l.Username)
		return nil, err.NoSuchUser
	}
	if !validPassword(u, l.Password) {
		log.Printf("Login password invalid? %s:%s\n", l.Username, l.Password)
		return nil, err.InconsistentState
	}
	log.Printf("Login loaded: %s:%s\n", l.Username, l.Password)
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
		log.Printf("No user named: %s\n", user)
		return err.NoSuchUser
	}
	if validPassword(u, pass) {
		writeLoginFile(user, pass)
		log.Printf("Login success")
		return err.OK
	}
	log.Printf("Invalid password for user: %s\n", user)
	return err.AuthenticateFail
}

// StoreUser to UserFile
func StoreUser(users entity.Users) {
	file, e := openFileRewrite(UserFile())
	err.LogFatalIfError(e)
	users.Serialize(file)
}

// StoreMeeting to MeetingFile
func StoreMeeting(meetings *entity.Meetings) {
	file, e := openFileRewrite(MeetingFile())
	err.LogFatalIfError(e)
	meetings.Serialize(file)
}
