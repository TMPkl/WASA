package migrations

import (
	"database/sql"
	"fmt"
)

func create_groups(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		photo_id INTEGER DEFAULT NULL,
		conversation_id INTEGER UNIQUE);

	CREATE TABLE IF NOT EXISTS groups_membership (
		group_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		PRIMARY KEY (group_id, username)
	);`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	return nil
}
