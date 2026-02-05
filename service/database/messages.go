package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database/attachments"
)

/*
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	conversation_id INTEGER NOT NULL,
	sender_username TEXT NOT NULL,
	content TEXT,
	timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
	attachment BLOB DEFAULT NULL,
	reaction text DEFAULT NULL,
	status TEXT CHECK( status IN ('sent','delivered','received','hidden') )
	CHECK (content IS NOT NULL OR attachment IS NOT NULL)
*/

type Message struct {
	ID                    int64
	ConversationID        uint
	SenderUsername        string
	Content               string
	Timestamp             time.Time
	Attachment            []byte // GOB encoded AttachmentsPack
	Reacted_to_message_id int64
	Replying_to_id        int64
	Status                string
}

func (db *appdbimpl) SaveMessage(username string, MessageContent string, ap attachments.AttachmentsPack, ConvID uint) (Message, error) {
	return db.SaveMessageWithReply(username, MessageContent, ap, ConvID, 0)
}

func (db *appdbimpl) SaveMessageWithReply(username string, MessageContent string, ap attachments.AttachmentsPack, ConvID uint, replyingToID int64) (Message, error) {
	exist, err := db.UserExists(username)
	if err != nil {
		return Message{}, fmt.Errorf("Database error: %w", err)
	}
	if !exist {
		return Message{}, errors.New("User does not exist")
	}
	message := Message{
		SenderUsername: username,
		Content:        MessageContent,
		Timestamp:      time.Now(),
		Status:         "sent",
		ConversationID: ConvID,
		Replying_to_id: replyingToID,
	}
	var attachmentsData []byte
	if ap.IsEmpty() {
		attachmentsData = nil
	} else {
		var err error
		attachmentsData, err = ap.ConvertToGOB()
		if err != nil {
			return Message{}, fmt.Errorf("Failed to convert attachments to GOB: %w", err)
		}
	}
	message.Attachment = attachmentsData

	res, err := db.c.Exec("INSERT INTO Messages (conversation_id, sender_username, content, timestamp, attachment, replying_to_id, status) VALUES (?, ?, ?, ?, ?, ?, ?)",
		message.ConversationID, message.SenderUsername, message.Content, message.Timestamp, message.Attachment, message.Replying_to_id, message.Status)
	if err != nil {
		return Message{}, fmt.Errorf("Failed to save message: %w", err)
	}
	resultID, err := res.LastInsertId()
	if err != nil {
		return Message{}, fmt.Errorf("Failed to retrieve last insert ID: %w", err)
	}
	message.ID = resultID
	return message, nil
}

func (db *appdbimpl) DoesUsersOwnConversation(unA string, unB string) (uint, error) {
	/*
		CREATE TABLE IF NOT EXISTS Private_conversations_memberships (
			conversation_id INTEGER NOT NULL,
			member_username TEXT NOT NULL,
			PRIMARY KEY (conversation_id, member_username));
	*/
	var conversationID uint
	err := db.c.QueryRow(`
		SELECT pcm1.conversation_id
		FROM Private_conversations_memberships pcm1
		JOIN Private_conversations_memberships pcm2
		  ON pcm1.conversation_id = pcm2.conversation_id
		WHERE pcm1.member_username = ? AND pcm2.member_username = ?
		LIMIT 1
	`, unA, unB).Scan(&conversationID)

	if err != nil {
		return 0, fmt.Errorf("Conversation not found: %w", err)
	}

	return conversationID, nil
}

func (db *appdbimpl) CreatePrivateConversation(username1 string, username2 string) (uint, error) {
	//  Check if both users exist
	exist1, err := db.UserExists(username1)
	if err != nil {
		return 0, fmt.Errorf("Database error: %w", err)
	}
	if !exist1 {
		return 0, errors.New("User " + username1 + " does not exist")
	}

	exist2, err := db.UserExists(username2)
	if err != nil {
		return 0, fmt.Errorf("Database error: %w", err)
	}
	if !exist2 {
		return 0, errors.New("User " + username2 + " does not exist")
	}

	// Create conversation
	res, err := db.c.Exec("INSERT INTO Conversations (type) VALUES (?)", "private")
	if err != nil {
		return 0, fmt.Errorf("Failed to create conversation: %w", err)
	}

	convID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Failed to retrieve conversation ID: %w", err)
	}

	// Add both users to the conversation
	_, err = db.c.Exec("INSERT INTO Private_conversations_memberships (conversation_id, member_username) VALUES (?, ?)",
		convID, username1)
	if err != nil {
		return 0, fmt.Errorf("Failed to add user1 to conversation: %w", err)
	}

	_, err = db.c.Exec("INSERT INTO Private_conversations_memberships (conversation_id, member_username) VALUES (?, ?)",
		convID, username2)
	if err != nil {
		return 0, fmt.Errorf("Failed to add user2 to conversation: %w", err)
	}

	return uint(convID), nil
}
func (db *appdbimpl) DeleteMessage(messageID string) error {
	_, err := db.c.Exec("DELETE FROM Messages WHERE id = ?", messageID)
	if err != nil {
		return fmt.Errorf("Failed to delete message: %w", err)
	}
	return nil
}
func (db *appdbimpl) MessageOwner(messageID string) (string, error) {
	var ownerUsername string
	err := db.c.QueryRow("SELECT sender_username FROM Messages WHERE id = ?", messageID).Scan(&ownerUsername)
	if err != nil {
		return "", fmt.Errorf("Failed bla bla bla cos tam cos : %w", err)
	}
	return ownerUsername, nil
}
func (db *appdbimpl) GetMessageByID(messageID string) (*Message, error) {
	var message Message
	err := db.c.QueryRow("SELECT id, conversation_id, sender_username, content, timestamp, COALESCE(attachment, ''), COALESCE(reacted_to_message_id, 0), COALESCE(replying_to_id, 0), status FROM Messages WHERE id = ?", messageID).
		Scan(&message.ID, &message.ConversationID, &message.SenderUsername, &message.Content, &message.Timestamp, &message.Attachment, &message.Reacted_to_message_id, &message.Replying_to_id, &message.Status)
	if err != nil {
		return nil, fmt.Errorf("Failed to bla bla bla fiu fiu fiu: %w", err)
	}
	return &message, nil
}

func (db *appdbimpl) UpdateMessageStatus(messageID string, newStatus string) error {
	validStatuses := map[string]bool{
		"sent":      true,
		"delivered": true,
		"received":  true,
		"hidden":    true,
	}

	if !validStatuses[newStatus] {
		return fmt.Errorf("invalid status: %s", newStatus)
	}

	result, err := db.c.Exec("UPDATE Messages SET status = ? WHERE id = ?", newStatus, messageID)
	if err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

func (db *appdbimpl) UserInConversation(username string, conversationID uint) (bool, error) {
	// Check if user is in private conversation
	var exists int
	err := db.c.QueryRow(`
		SELECT 1
		FROM Private_conversations_memberships
		WHERE member_username = ? AND conversation_id = ?
		LIMIT 1
	`, username, conversationID).Scan(&exists)

	if err == nil {
		return true, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("Database error checking private conversation: %w", err)
	}

	// Check if user is in group conversation
	err = db.c.QueryRow(`
		SELECT 1
		FROM Groups g
		JOIN Groups_memberships gm ON g.id = gm.group_id
		WHERE gm.member_username = ? AND g.conversation_id = ?
		LIMIT 1
	`, username, conversationID).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("Database error checking group conversation: %w", err)
	}

	return true, nil
}
func (db *appdbimpl) ConversationIDfromMessageID(messageID string) (uint, error) {
	var convID uint
	err := db.c.QueryRow("SELECT conversation_id FROM Messages WHERE id = ?", messageID).Scan(&convID)
	if err != nil {
		return 0, fmt.Errorf("Failed to get conversation ID from message ID: %w", err)
	}
	return convID, nil
}

func (db *appdbimpl) ReactToMessage(reactedToMessagID int, reaction string, username string, conversationID uint) (uint, error) {
	exist, err := db.UserExists(username)
	if err != nil {
		return 0, fmt.Errorf("Database error: %w", err)
	}
	if !exist {
		return 0, errors.New("User does not exist")
	}
	//dodaj sprawdzenie czy user juz nie zareagowal na ta wiadomosc TODO
	/*
		type Message struct {
		ID             int64
		ConversationID uint
		SenderUsername string
		Content        string
		Timestamp      time.Time
		Attachment     []byte // GOB encoded AttachmentsPack
		Reaction       string
		Status         string
	*/
	message := Message{
		SenderUsername:        username,
		Content:               reaction,
		Timestamp:             time.Now(),
		Status:                "sent",
		ConversationID:        conversationID,
		Reacted_to_message_id: int64(reactedToMessagID),
	}

	_, err = db.c.Exec("INSERT INTO Messages (conversation_id, sender_username, content, timestamp, reacted_to_message_id, status) VALUES (?, ?, ?, ?, ?, ?)",
		message.ConversationID, message.SenderUsername, message.Content, message.Timestamp, message.Reacted_to_message_id, message.Status)
	if err != nil {
		return 0, fmt.Errorf("Failed to save reaction message: %w", err)
	}
	return uint(message.ID), nil

}
func (db *appdbimpl) GetReactionsID(baseMessageID int) ([]int, bool) {
	res, err := db.c.Query("Select r.id from Messages b join Messages r on b.id = r.reacted_to_message_id where b.id = ?", baseMessageID)
	if err != nil {
		return nil, false
	}
	defer res.Close()
	var reactionsIDs []int
	for res.Next() {
		var reactionID int
		err := res.Scan(&reactionID)
		if err != nil {
			return nil, false
		}
		reactionsIDs = append(reactionsIDs, reactionID)
	}
	if err := res.Err(); err != nil {
		return nil, false
	}
	return reactionsIDs, true
}

func (db *appdbimpl) GetReactionIDByUsernameAndBaseMessageID(username string, baseMessageID int) (int, error) {
	var reactionMessageID int
	err := db.c.QueryRow("Select r.id from Messages b join Messages r on b.id = r.reacted_to_message_id where b.id = ? and r.sender_username = ?", baseMessageID, username).Scan(&reactionMessageID)
	if err != nil {
		return 0, fmt.Errorf("Failed to get reaction message ID: %w", err)
	}
	return reactionMessageID, nil
}

// GetConversationHistory retrieves the message history for a conversation
// Returns messages ordered by timestamp (newest first) with a limit
func (db *appdbimpl) GetConversationHistory(conversationID uint, limit int) ([]Message, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `
		SELECT id, conversation_id, sender_username, content, timestamp, 
		       COALESCE(attachment, ''), COALESCE(reacted_to_message_id, 0), COALESCE(replying_to_id, 0), status
		FROM Messages
		WHERE conversation_id = ? AND reacted_to_message_id IS NULL
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := db.c.Query(query, conversationID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.SenderUsername,
			&msg.Content,
			&msg.Timestamp,
			&msg.Attachment,
			&msg.Reacted_to_message_id,
			&msg.Replying_to_id,
			&msg.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	// Reverse to get oldest first (chronological order)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
func (db *appdbimpl) GetAttachmentFromMessage(messageID int) ([]byte, error) {
	var attachmentData []byte
	err := db.c.QueryRow("SELECT attachment FROM Messages WHERE id = ?", messageID).Scan(&attachmentData)
	if err != nil {
		return nil, fmt.Errorf("Failed to get attachment from message: %w", err)
	}
	return attachmentData, nil
}
