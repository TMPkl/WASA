package migrations

import (
	"database/sql"
	"fmt"
)

func create_messages(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS Messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		conversation_id INTEGER NOT NULL,
		sender_username TEXT NOT NULL,
		content TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		attachment BLOB DEFAULT NULL,
		reaction text DEFAULT NULL,
		status TEXT CHECK( status IN ('sent','delivered','received','hidden') ), 
		CHECK (content IS NOT NULL OR attachment IS NOT NULL)
	);`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	return nil
}
