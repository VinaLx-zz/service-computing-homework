package sqldao

import (
	"args"
	"database/sql"
	"log"
	"time"
	"user"
)

type userDB struct {
	db *sql.DB
}

var d *userDB

func prepareDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", *args.DB)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS User (
			ID INTEGER NOT NULL PRIMARY KEY,
			Username TEXT NOT NULL,
			Password TEXT NOT NULL,
			SignUpDate TEXT
		);`)
	if err != nil {
		return nil, err
	}
	log.Print("database initialized")
	return db, err
}

func sqliteDateToTime(date string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05.000000-07:00", date)
	if err != nil {
		panic(err)
	}
	return t
}

// Get dao
func Get() (user.Dao, error) {
	if d == nil {
		db, err := prepareDB()
		if err != nil {
			return nil, err
		}
		d = &userDB{db}
	}
	return d, nil
}

func (ud *userDB) StoreUser(u *user.User) error {
	_, err := ud.db.Exec(
		`INSERT INTO User (ID, Username, Password, SignUpDate)
		VALUES (?, ?, ?, ?)`, u.ID, u.Username, u.Password, u.SignUpDate)
	return err
}

func (ud *userDB) GetAllUsers() ([]*user.User, error) {
	row, err := ud.db.Query("SELECT * FROM User")
	if err != nil {
		return nil, err
	}
	result := make([]*user.User, 0)
	for row.Next() {
		var u user.User
		var date string
		row.Scan(&u.ID, &u.Username, &u.Password, &date)
		u.SignUpDate = sqliteDateToTime(date)
		result = append(result, &u)
	}
	return result, nil
}

func (ud *userDB) GetUser(uid uint64) (*user.User, error) {
	row, err := ud.db.Query("SELECT * FROM User WHERE uid = ?", uid)
	if err != nil {
		return nil, err
	}
	if row.Next() {
		var u user.User
		var date string
		row.Scan(&u.ID, &u.Username, &u.Password, &date)
		u.SignUpDate = sqliteDateToTime(date)
		return &u, nil
	}
	return nil, nil
}
