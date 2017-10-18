package cmd

import (
	"crypto/md5"
	"entity"
	"err"
	"fmt"
	"log"
	"model"
	"os"
	"strings"
	"time"
)

func printWrongLoginState(action string, required bool) int {
	var s string
	if required {
		s = "login"
	} else {
		s = "logout"
	}
	fmt.Fprintf(os.Stderr, "Action %s requires an %s state\n", action, s)
	return int(err.WrongLoginState)
}

func printMeetingDoesntExist(title string) int {
	fmt.Fprintf(os.Stderr, "meeting doesn't exist: %s\n", title)
	return int(err.NoSuchMeeting)
}

func printNotAHost(title string, user string) int {
	fmt.Fprintf(os.Stderr, "meeting '%s' is not hosted by '%s'", title, user)
	return int(err.NotEnoughPrivilege)
}

func printUserDoesntExist(username string) int {
	fmt.Fprintf(os.Stderr, "user doesn't exist: %s\n", username)
	return int(err.NoSuchUser)
}

func printInvalidTimeFormat(time string) int {
	fmt.Fprintf(os.Stderr, "invalid time format: %s\n", time)
	return int(err.InvalidTime)
}

func loadLogin(us entity.Users) *entity.User {
	u, e := model.LoadLogin(us)
	if e == err.OK {
		return u
	}
	log.Fatalf("something wrong with login file, error: %d", int(e))
	return nil
}

// Register a user
func Register(user, pass, mail, phone string) int {
	users := model.LoadUsers()
	passhash := fmt.Sprintf("%x", md5.Sum([]byte(pass)))
	log.Printf("password hash for '%s': %s\n", pass, passhash)
	if !users.Add(&entity.User{
		Username: user,
		Password: passhash,
		Mail:     mail,
		Phone:    phone,
	}) {
		fmt.Fprintf(os.Stderr, "there's another user with username %s\n", user)
		return int(err.DuplicateUser)
	}
	model.StoreUser(users)
	return 0
}

// Login Command
func Login(user, pass string) int {
	users := model.LoadUsers()
	if loadLogin(users) != nil {
		return printWrongLoginState("login", false)
	}
	if model.Login(users, user, pass) != err.OK {
		fmt.Fprintln(os.Stderr, "Authentication Fail")
		return int(err.AuthenticateFail)
	}
	return 0
}

// Logout Command
func Logout() int {
	if model.Logout() {
		fmt.Printf("logout success\n")
	} else {
		return printWrongLoginState("logout", true)
	}
	return 0
}

// ShowUsers print all users when logged in
func ShowUsers() int {
	users := model.LoadUsers()
	if loadLogin(users) == nil {
		return printWrongLoginState("ShowUser", true)
	}
	fmt.Println("Username Email Phone")
	for _, u := range users {
		fmt.Printf("'%s' '%s' '%s'\n", u.Username, u.Mail, u.Phone)
	}
	return 0
}

// DeleteUser delete current login user, and removed from its meeting
func DeleteUser() int {
	users := model.LoadUsers()
	u := loadLogin(users)
	if u == nil {
		return printWrongLoginState("DeleteUser", true)
	}
	if users.Remove(u) == nil {
		log.Fatalln("Login user cannot be removed !?")
	}
	meetings := model.LoadMeetings()
	meetings.RemoveAll(u)
	model.StoreMeeting(meetings)
	model.StoreUser(users)
	return 0
}

const timeLayout = "2006-01-02"

func hostMeeting(ms *entity.Meetings, m *entity.Meeting) int {
	e := ms.Host(m)
	switch e {
	case err.InvalidTime:
		fmt.Fprintln(os.Stderr, "meeting should end later than start")
	case err.DuplicateMeeting:
		fmt.Fprintf(os.Stderr, "there's another meeting with title: %s\n", m.Title)
	case err.TimeConflict:
		fmt.Fprintln(os.Stderr, "there are time conflict of some participants")
	case err.OK:
		model.StoreMeeting(ms)
		fmt.Println("meeting hosted")
	}
	return int(e)
}

// HostMeeting of specified parameter
func HostMeeting(title string, parts []string, start, end string) int {
	users := model.LoadUsers()
	host := loadLogin(users)
	if host == nil {
		return printWrongLoginState("HostMeeting", true)
	}
	s, es := time.Parse(timeLayout, start)
	if es != nil {
		return printInvalidTimeFormat(start)
	}
	e, ee := time.Parse(timeLayout, end)
	if ee != nil {
		return printInvalidTimeFormat(end)
	}
	participants := entity.NewUsers()
	for _, p := range parts {
		u := users.Lookup(p)
		if u == nil {
			return printUserDoesntExist(p)
		}
		participants.Add(u)
	}
	meeting := &entity.Meeting{
		Title:        title,
		Host:         host,
		Participants: participants,
		Start:        s,
		End:          e,
	}
	meetings := model.LoadMeetings()
	return hostMeeting(meetings, meeting)
}

// CancelMeeting with specific title hosted by logined user
func CancelMeeting(title string) int {
	users := model.LoadUsers()
	host := loadLogin(users)
	if host == nil {
		return printWrongLoginState("CancelMeeting", true)
	}
	meetings := model.LoadMeetings()
	meeting := meetings.Lookup(title)
	if meeting == nil {
		return printMeetingDoesntExist(title)
	}
	if meeting.Host.Username != host.Username {
		return printNotAHost(title, host.Username)
	}
	meetings.Cancel(title)
	model.StoreMeeting(meetings)
	return 0
}

// QuitMeeting of specific title of the login user
func QuitMeeting(title string) int {
	users := model.LoadUsers()
	meetings := model.LoadMeetings()
	user := loadLogin(users)
	e := meetings.Remove(title, user)
	switch e {
	case err.NoSuchMeeting:
		printMeetingDoesntExist(title)
	case err.NoSuchUser:
		fmt.Fprintf(
			os.Stderr, "user '%s' is not a participant of meeting '%s'\n",
			user.Username, title)
	case err.OK:
		model.StoreMeeting(meetings)
	}
	return int(e)
}

// AddParticipant add user to the hosted meeting
func AddParticipant(title string, username string) int {
	users := model.LoadUsers()
	meetings := model.LoadMeetings()
	host := loadLogin(users)
	if host == nil {
		return printWrongLoginState("AddParticipant", true)
	}
	meeting := meetings.Lookup(title)
	if meeting == nil {
		return printMeetingDoesntExist(title)
	}
	if meeting.Host.Username != host.Username {
		return printNotAHost(title, host.Username)
	}
	addedUser := users.Lookup(username)
	if addedUser == nil {
		return printUserDoesntExist(username)
	}
	e := meetings.Add(title, addedUser)
	switch e {
	case err.NoSuchMeeting:
		log.Fatalln("meeting should be checked before??")
	case err.DuplicateUser:
		fmt.Fprintf(
			os.Stderr, "user '%s' is already a participant of meeting '%s'",
			username, title)
	case err.TimeConflict:
		fmt.Fprintf(os.Stderr,
			"there is time conflict of user '%s' and meeting '%s'",
			username, title)
	case err.OK:
		model.StoreMeeting(meetings)
	default:
		log.Fatalf("unexpected error: %d\n", int(e))
	}
	return int(e)
}

// RemoveParticipant from a specific hosted meeting
func RemoveParticipant(title string, username string) int {
	users := model.LoadUsers()
	host := loadLogin(users)
	if host == nil {
		return printWrongLoginState("RemoveParticipant", true)
	}
	meetings := model.LoadMeetings()
	meeting := meetings.Lookup(title)
	if meeting == nil {
		return printMeetingDoesntExist(title)
	}
	if meeting.Host.Username != host.Username {
		return printNotAHost(title, host.Username)
	}
	user := users.Lookup(username)
	if user == nil {
		return printUserDoesntExist(username)
	}
	e := meetings.Remove(title, user)
	if e != err.OK {
		log.Fatalf("error should be all checked! %d", e)
	}
	model.StoreMeeting(meetings)
	return 0
}

// ClearMeetings of the login user
func ClearMeetings() int {
	users := model.LoadUsers()
	host := loadLogin(users)
	if host == nil {
		return printWrongLoginState("ClearMeeting", true)
	}
	meetings := model.LoadMeetings()
	meetings.CancelAll(host)
	model.StoreMeeting(meetings)
	return 0
}

func printMeeting(m *entity.Meeting) {
	parts := make([]string, 0)
	for _, u := range m.Participants {
		parts = append(parts, u.Username)
	}
	fmt.Printf("title: %s\n\thost: %s\n\ttime: %s to %s\n\tparticipants: %s\n",
		m.Title, m.Host.Username,
		m.Start.Format(timeLayout), m.End.Format(timeLayout),
		strings.Join(parts, ", "))
}

// QueryMeeting overlapped with specific time interval
func QueryMeeting(start, end string) int {
	users := model.LoadUsers()
	user := loadLogin(users)
	if user == nil {
		return printWrongLoginState("QueryMeeting", true)
	}
	s, es := time.Parse(timeLayout, start)
	if es != nil {
		return printInvalidTimeFormat(start)
	}
	e, ee := time.Parse(timeLayout, end)
	if ee != nil {
		return printInvalidTimeFormat(end)
	}
	if e.Before(s) {
		fmt.Fprintf(os.Stderr, "invalid interval %s - %s", start, end)
		return int(err.InvalidTime)
	}
	meetings := model.LoadMeetings()
	for _, m := range meetings.Related(user.Username) {
		if e.After(m.Start) && s.Before(m.End) {
			printMeeting(m)
		}
	}
	return 0
}
