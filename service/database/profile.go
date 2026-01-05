package database

import "errors"

func (db *appdbimpl) UserExists(username string) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=?)", username).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) AddNewUser(username string) error {
	exist, err := db.UserExists(username)
	if err != nil {
		return errors.New("Database error")
	}
	if exist {
		return errors.New("User already exists")
	}
	_, err = db.c.Exec("INSERT INTO users (username) VALUES (?)", username)
	return err
}
