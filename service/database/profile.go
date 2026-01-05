package database

import (
	"database/sql"
	"errors"
)

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

func (db *appdbimpl) AddProfilePhoto(username string, photoData []byte) error {
	//1. uzytkownik istnieje w bazie danych
	//2. wstaw zdjecie do zdjec
	//3. wez id zdjecia
	//4. podepnij ip do profilu
	exist, err := db.UserExists(username)
	if err != nil {
		return errors.New("Database error")
	}
	if !exist {
		return errors.New("User does not exist")
	}
	res, err := db.c.Exec("INSERT INTO Users_photos(photo_data) VALUES(?)", photoData)

	if err != nil {
		return errors.New("Database error")
	}
	new_photo_id, _ := res.LastInsertId()

	//sprawdz czy uzytkownik ma zdjecie
	//A jesli ma to wywal tamto z bazy, dopisz mu ID nowego
	//B jeśli nie to wpisz mu ID nowego
	var old_photo_id sql.NullInt64
	err = db.c.QueryRow("SELECT photo_id FROM users WHERE username=?)", username).Scan(&old_photo_id)
	if err != nil {
		return errors.New("Database error")
	}
	if old_photo_id.Valid {
		// update starego zdjecia
		_, err = db.c.Exec("DELETE FROM Users_photos WHERE photo_id=?", old_photo_id.Int64)
		///error do przemyslenia co z nim zrobic
	}

	_, err = db.c.Exec("Update Users SET photo_id=? where username=?", new_photo_id, username) //zapisz nowe id_zdjecia dla usera

	return nil
}
