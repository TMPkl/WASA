package migrations

import (
	"database/sql"
	"fmt"
)

func create_user_photos(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users_photos (
		photo_id INTEGER PRIMARY KEY AUTOINCREMENT,
		photo_data BLOB NOT NULL
	);`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	return nil
}
