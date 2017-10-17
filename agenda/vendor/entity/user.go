package entity

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
