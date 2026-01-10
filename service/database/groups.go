package database

import (
	"fmt"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database/attachments"
)

func (db *appdbimpl) CreateGroup(groupName string, firstMember string) (uint, error) {
	// Start a transaction
	tx, err := db.c.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() // rollback if not committed

	//1crate conversation and save its id
	//2create group based on conversation id
	//3 add members to a conversation by inserting into Groups_memberships

	//1 Create a new group conversation
	res, err := tx.Exec("INSERT INTO Conversations (type, photo_id) VALUES ('group', NULL)")
	if err != nil {
		return 0, fmt.Errorf("failed to create group conversation: %w", err)
	}
	conversationID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve conversation ID: %w", err)
	}

	//2 creete the group
	res, err = tx.Exec("INSERT INTO Groups (name, conversation_id) VALUES (?, ?)", groupName, conversationID)
	if err != nil {
		return 0, fmt.Errorf("failed to create group: %w", err)
	}
	groupID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve group ID: %w", err)
	}

	//3  adding members
	_, err = tx.Exec("INSERT INTO Groups_memberships (group_id, member_username) VALUES (?, ?)", groupID, firstMember)
	if err != nil {
		return 0, fmt.Errorf("failed to add member %s to group: %w", firstMember, err)
	}
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return uint(groupID), nil
}

func (db *appdbimpl) GetGroupMembers(groupID uint) ([]string, error) {
	rows, err := db.c.Query("SELECT member_username FROM Groups_memberships WHERE group_id = ?", groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to query group members: %w", err)
	}
	defer rows.Close()

	var members []string
	for rows.Next() {
		var member string
		if err := rows.Scan(&member); err != nil {
			return nil, fmt.Errorf("failed to scan member username: %w", err)
		}
		members = append(members, member)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over group members: %w", err)
	}

	return members, nil
}

func (db *appdbimpl) AddMemberToGroup(groupID uint, username string) error {
	_, err := db.c.Exec("INSERT INTO Groups_memberships (group_id, member_username) VALUES (?, ?)", groupID, username)
	if err != nil {
		return fmt.Errorf("failed to add member to group: %w", err)
	}
	return nil
}

func (db *appdbimpl) RemoveMemberFromGroup(groupID uint, username string) error {
	_, err := db.c.Exec("DELETE FROM Groups_memberships WHERE group_id = ? AND member_username = ?", groupID, username)
	if err != nil {
		return fmt.Errorf("failed to remove member from group: %w", err)
	}
	return nil
}

func (db *appdbimpl) IsUserInGroup(groupID uint, username string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Groups_memberships WHERE group_id = ? AND member_username = ?", groupID, username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user in group: %w", err)
	}
	return count > 0, nil
}
func (db *appdbimpl) SaveGroupMessage(username string, groupID uint, messageContent string, ap attachments.AttachmentsPack) (Message, error) {
	// glownie po to to zrobilem aby wysylac wiadomosc na grupe na podstawie groupID a nie conversationID
	var conversationID uint
	err := db.c.QueryRow("SELECT conversation_id FROM Groups WHERE id = ?", groupID).Scan(&conversationID)
	if err != nil {
		return Message{}, fmt.Errorf("failed to get conversation ID for group: %w", err)
	}

	message, err := db.SaveMessage(username, messageContent, ap, conversationID)
	if err != nil {
		return Message{}, fmt.Errorf("failed to save group message: %w", err)
	}

	return message, nil
}
func (db *appdbimpl) RenameGroup(groupID uint, newName string) error {
	_, err := db.c.Exec("UPDATE Groups SET name = ? WHERE id = ?", newName, groupID)
	if err != nil {
		return fmt.Errorf("failed to rename group: %w", err)
	}
	return nil
}
func (db *appdbimpl) UpdateGroupPhoto(groupID uint, photoData []byte) error {

	tx, err := db.c.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() //deferrrrr

	res, err := tx.Exec("INSERT INTO Users_photos(photo_data) VALUES(?)", photoData)
	if err != nil {
		return fmt.Errorf("failed to insert photo: %w", err)
	}

	newPhotoID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get new photo ID: %w", err)
	}

	_, err = tx.Exec("UPDATE Groups SET photo_id = ? WHERE id = ?", newPhotoID, groupID)
	if err != nil {
		return fmt.Errorf("failed to update group with new photo ID: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
func (db *appdbimpl) GroupExists(groupID uint) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Groups WHERE id = ?", groupID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if group exists: %w", err)
	}
	return count > 0, nil
}
