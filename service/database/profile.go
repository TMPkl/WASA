package database

import (
	"database/sql"
	"errors"
	"fmt"
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
	exist, err := db.UserExists(username)
	if err != nil {
		return fmt.Errorf("Database error: %w", err)
	}
	if !exist {
		return errors.New("User does not exist")
	}

	tx, err := db.c.Begin()
	if err != nil {
		return fmt.Errorf("Failed to start transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Wstaw nowy obraz
	res, err := tx.Exec("INSERT INTO Users_photos(photo_data) VALUES(?)", photoData)
	if err != nil {
		return fmt.Errorf("Failed to insert photo: %w", err)
	}

	newPhotoID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("Failed to get new photo ID: %w", err)
	}

	//stare photo_id
	var oldPhotoID sql.NullInt64
	err = tx.QueryRow("SELECT photo_id FROM users WHERE username=?", username).Scan(&oldPhotoID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("Failed to query old photo_id: %w", err)
	}

	// Usuń stare zdjęcie
	if oldPhotoID.Valid {
		_, err = tx.Exec("DELETE FROM Users_photos WHERE photo_id=?", oldPhotoID.Int64)
		if err != nil {
			return fmt.Errorf("Failed to delete old photo: %w", err)
		}
	}

	//se new photo id in user
	_, err = tx.Exec("UPDATE Users SET photo_id=? WHERE username=?", newPhotoID, username)
	if err != nil {
		return fmt.Errorf("Failed to update user: %w", err)
	}

	return tx.Commit() // jeśli wszystko OK
}

func (db *appdbimpl) GetProfilePhoto(username string) ([]byte, error) {
	var photoData []byte
	err := db.c.QueryRow(`
		SELECT up.photo_data
		FROM Users u
		JOIN Users_photos up ON u.photo_id = up.photo_id
		WHERE u.username = ?`, username).Scan(&photoData)
	if err != nil {
		return nil, fmt.Errorf("Failed to get profile photo: %w", err)
	}
	return photoData, nil
}
