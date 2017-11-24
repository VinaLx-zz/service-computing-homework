package ormdao

import (
	"args"
	"log"
	"user"

	"github.com/jinzhu/gorm"
)

type orm struct {
	db *gorm.DB
}

var o *orm

func prepareDB(createTable bool) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", *args.DB)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return nil, err
	}
	db.SingularTable(true)
	db.LogMode(*args.Log)
	if createTable {
		db = db.CreateTable(&user.User{})
		if db.Error != nil {
			db.Close()
			db, err = prepareDB(false)
		}
	}
	return db, nil
}

// Get dao
func Get() (user.Dao, error) {
	if o == nil {
		db, err := prepareDB(true)
		if err != nil {
			return nil, err
		}
		o = &orm{db}
	}
	return o, nil
}

func (o *orm) StoreUser(u *user.User) error {
	if o.db.Create(u).Error != nil {
		return o.db.Error
	}
	return nil
}

func (o *orm) GetAllUsers() ([]*user.User, error) {
	users := make([]*user.User, 0)
	if o.db.Find(&users).Error != nil {
		return nil, o.db.Error
	}
	return users, nil
}

func (o *orm) GetUser(uid uint64) (*user.User, error) {
	var u user.User
	if o.db.Find(&u, uid).Error != nil {
		return nil, o.db.Error
	}
	if u.Username != "" {
		return &u, nil
	}
	return nil, nil
}
