package entity

import (
	"encoding/json"
	"io"
)

// User struct
// Password should be the md5 hash of password clear text
type User struct {
	Username string
	Password string
	Mail     string
	Phone    string
}

// Users is the essentially a map of username to pointer of User,
// which disallows duplication of username
type Users map[string]*User

// NewUsers create an empty user set
func NewUsers() Users {
	return make(Users)
}

// Lookup returns the pointer to the user of specific username
// nil if not found
func (users Users) Lookup(username string) *User {
	return users[username]
}

// Add returns false if there's already a user with the same username in map
// returns true if the add success
func (users Users) Add(user *User) bool {
	if user == nil {
		return false
	}
	u := users.Lookup(user.Username)
	if u != nil {
		return false
	}
	users[user.Username] = user
	return true
}

// Remove returns the removed user in Users, nil if nothing removed
func (users Users) Remove(user *User) *User {
	if user == nil {
		return nil
	}
	u := users.Lookup(user.Username)
	if u == nil {
		return nil
	}
	delete(users, user.Username)
	return u
}

// Slice returns the slice of the user in the map
func (users Users) Slice() []*User {
	s := make([]*User, 0, len(users))
	for _, u := range users {
		s = append(s, u)
	}
	return s
}

// Serialize users to specific writer
func (users Users) Serialize(w io.Writer) {
	encoder := json.NewEncoder(w)
	for _, u := range users {
		encoder.Encode(u)
	}
}

// DeserializeUser restore Users from the reader in the format of serialization
func DeserializeUser(r io.Reader) (Users, error) {
	decoder := json.NewDecoder(r)
	users := make(Users)
	for {
		u := new(User)
		if err := decoder.Decode(u); err == io.EOF {
			return users, nil
		} else if err != nil {
			return nil, err
		}
		users.Add(u)
	}
}
