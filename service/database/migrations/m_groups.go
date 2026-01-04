package migrations

import (
	"database/sql"
	"fmt"
)

func create_groups(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS Groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		conversation_id INTEGER UNIQUE,
		photo_id INTEGER DEFAULT NULL
	);`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	return nil
}
