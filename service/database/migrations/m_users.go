package migrations

import (
	"database/sql"
	"fmt"
)

func create_users(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT NOT NULL PRIMARY KEY,
		bearerToken TEXT,
		photo_id INTEGER DEFAULT NULL
	);`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	return nil
}
