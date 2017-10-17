package entity

import (
	"err"
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
	return e1.After(s1) && s1.Before(e2)
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
}

// Host create a meeting in meetings if all constraints are satisfied:
// title shouldn't be duplicate
// all user including host should have time
func (ms *Meetings) Host(
	title string, host *User, other Users, start, end time.Time) err.Err {
	if !end.After(start) {
		return err.InvalidTime
	}
	if ms.Has(title) {
		return err.DuplicateMeeting
	}
	other.Add(host)
	for _, u := range other {
		for _, m := range ms.Related(u.Username) {
			if overlapped(m.Start, m.End, start, end) {
				return err.InvalidTime
			}
		}
	}
	ms.host(&Meeting{
		Title:        title,
		Host:         host,
		Participants: other,
		Start:        start,
		End:          end,
	})
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

// Remove remove user from specific meeting and returns OK if success
// NoSuchUser or NoSuchMeeting when error
func (ms *Meetings) Remove(title string, user *User) err.Err {
	m := ms.meetings[title]
	if m == nil {
		return err.NoSuchMeeting
	}
	if m.Participants.Remove(user) == nil {
		return err.NoSuchUser
	}
	if user.Username != m.Host.Username {
		delete(ms.relation[user.Username], m.Title)
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
