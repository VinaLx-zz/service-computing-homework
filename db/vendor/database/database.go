package database

import "user"

// UserDao is User Data Access Object
type UserDao interface {
	StoreUser(u *user.User) error
	GetAllUsers() ([]*user.User, error)
	GetUser(uid uint64) (*user.User, error)
}

// PureSQLDao returns the Dao with database/sql behind it
func PureSQLDao() UserDao {
	return nil
}

// ORMDao returns the Dao with ORM behind it
func ORMDao() UserDao {
	return nil
}
