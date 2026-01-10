package migrations

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) error {
	if err := create_users(db); err != nil {
		return fmt.Errorf("migration create_user failed: %w", err)
	}
	if err := create_user_photos(db); err != nil {
		return fmt.Errorf("migration create_user_photos failed: %w", err)
	}
	if err := create_groups(db); err != nil {
		return fmt.Errorf("migration create_groups failed: %w", err)
	}
	if err := create_conversations(db); err != nil {
		return fmt.Errorf("migration create_conversations failed: %w", err)
	}
	if err := create_messages(db); err != nil {
		return fmt.Errorf("migration create_messages failed: %w", err)
	}
	return nil
}
