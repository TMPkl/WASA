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

func (db *appdbimpl) UpdateUsername(oldUsername, newUsername string) error {
	exist, err := db.UserExists(newUsername)
	if err != nil {
		return errors.New("Database error")
	}
	if exist {
		return errors.New("New username already exists")
	}
	_, err = db.c.Exec("UPDATE users SET username=? WHERE username=?", newUsername, oldUsername)
	return err
}
