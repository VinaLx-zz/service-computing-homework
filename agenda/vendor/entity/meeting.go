package entity

import (
	"encoding/json"
	"err"
	"io"
	"log"
	"time"
)

// Meeting struct
type Meeting struct {
	Title        string
	Host         *User
	Participants Users
	Start        time.Time
	End          time.Time
}

// Meetings can hold meetings that logically existable
type Meetings struct {
	meetings map[string]*Meeting
	relation map[string]map[string]*Meeting
}

// NewMeetings returns an empty meeting set
func NewMeetings() *Meetings {
	return &Meetings{
		meetings: make(map[string]*Meeting),
		relation: make(map[string]map[string]*Meeting),
	}
}

// Lookup find the meeting with specific title
func (ms *Meetings) Lookup(title string) *Meeting {
	return ms.meetings[title]
}

// Slice returns the array of all meetings in Meetings
func (ms *Meetings) Slice() []*Meeting {
	a := make([]*Meeting, 0, len(ms.meetings))
	for _, m := range ms.meetings {
		a = append(a, m)
	}
	return a
}

// Has tests if there's a meeting with specific title
func (ms *Meetings) Has(title string) bool {
	return ms.meetings[title] != nil
}

// Related returns all meetings that are related to some user
// (participated or hosted)
func (ms *Meetings) Related(user string) map[string]*Meeting {
	related := ms.relation[user]
	if related == nil {
		return make(map[string]*Meeting)
	}
	return related
}

func overlapped(s1, e1, s2, e2 time.Time) bool {
	return e1.After(s2) && s1.Before(e2)
}

func (ms *Meetings) addRelatedMeeting(u *User, m *Meeting) {
	meetings := ms.relation[u.Username]
	if meetings == nil {
		ms.relation[u.Username] = map[string]*Meeting{m.Title: m}
	} else {
		meetings[m.Title] = m
	}
}

func (ms *Meetings) host(m *Meeting) {
	ms.meetings[m.Title] = m
	for _, u := range m.Participants {
		ms.addRelatedMeeting(u, m)
	}
	ms.addRelatedMeeting(m.Host, m)
}

// Host create a meeting in meetings if all constraints are satisfied:
// title shouldn't be duplicate
// all user including host should have time
func (ms *Meetings) Host(m *Meeting) err.Err {
	if !m.End.After(m.Start) {
		return err.InvalidTime
	}
	if ms.Has(m.Title) {
		return err.DuplicateMeeting
	}
	for _, u := range m.Participants {
		for _, um := range ms.Related(u.Username) {
			if overlapped(m.Start, m.End, um.Start, um.End) {
				log.Printf(
					"time conflict for '%s': %s to %s\n", u.Username,
					um.Start.Format("2006-01-02"), um.End.Format("2006-01-02"))
				return err.TimeConflict
			}
		}
	}
	ms.host(m)
	return err.OK
}

func (ms *Meetings) cancel(m *Meeting) {
	for _, u := range m.Participants {
		delete(ms.relation[u.Username], m.Title)
	}
	delete(ms.meetings, m.Title)
}

// Cancel directly remove the meeting of specific title from Meetings
func (ms *Meetings) Cancel(title string) err.Err {
	m := ms.meetings[title]
	if m == nil {
		return err.NoSuchMeeting
	}
	ms.cancel(m)
	return err.OK
}

// CancelAll meetings hosted by specific user
func (ms *Meetings) CancelAll(host *User) {
	for _, m := range ms.Related(host.Username) {
		if m.Host.Username == host.Username {
			ms.cancel(m)
		}
	}
}

func (ms *Meetings) remove(m *Meeting, user *User) {
	m.Participants.Remove(user)
	if user.Username != m.Host.Username || m.Participants.Size() == 0 {
		delete(ms.relation[user.Username], m.Title)
	}
	if m.Participants.Size() == 0 {
		delete(ms.meetings, m.Title)
	}
}

// Remove remove user from specific meeting and returns OK if success
// NoSuchUser or NoSuchMeeting when error
func (ms *Meetings) Remove(title string, user *User) err.Err {
	m := ms.meetings[title]
	if m == nil {
		return err.NoSuchMeeting
	}
	if m.Participants.Lookup(user.Username) == nil {
		return err.NoSuchUser
	}
	ms.remove(m, user)
	return err.OK
}

// RemoveAll removes the user from all the meeting
func (ms *Meetings) RemoveAll(user *User) err.Err {
	ms.CancelAll(user)
	related := make([]*Meeting, 0)
	for _, m := range ms.Related(user.Username) {
		related = append(related, m)
	}
	for _, m := range related {
		ms.remove(m, user)
	}
	return err.OK
}

// Add add specific user to specific meeting
func (ms *Meetings) Add(title string, user *User) err.Err {
	m := ms.meetings[title]
	if m == nil {
		return err.NoSuchMeeting
	}
	if m.Participants.Lookup(user.Username) != nil {
		return err.DuplicateUser
	}
	for _, um := range ms.Related(user.Username) {
		if overlapped(um.Start, um.End, m.Start, m.End) {
			return err.TimeConflict
		}
	}
	m.Participants.Add(user)
	ms.addRelatedMeeting(user, m)
	return err.OK
}

// Serialize meetings to specific writer
func (ms *Meetings) Serialize(w io.Writer) {
	encoder := json.NewEncoder(w)
	for _, m := range ms.meetings {
		encoder.Encode(m)
	}
}

// DeserializeMeeting restore meetings from the reader in the format of
// the serialization
func DeserializeMeeting(r io.Reader) (*Meetings, error) {
	decoder := json.NewDecoder(r)
	ms := NewMeetings()
	for {
		m := new(Meeting)
		if err := decoder.Decode(m); err == io.EOF {
			return ms, nil
		} else if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		ms.host(m)
	}
}
