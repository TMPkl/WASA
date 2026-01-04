package migrations

import (
	"database/sql"
	"fmt"
)

func create_conversations(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS Conversations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT CHECK( type IN ('private','group') ) NOT NULL,
		photo_id INTEGER,
			CHECK (
			(type = 'group' and photo_id IS NOT NULL)
			OR
			(type = 'private' and photo_id IS NULL))
	);
	CREATE TABLE IF NOT EXISTS Private_conversations_memberships (
		conversation_id INTEGER NOT NULL,
		member_username TEXT NOT NULL,
		PRIMARY KEY (conversation_id, member_username));
		
		`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	return nil
}
