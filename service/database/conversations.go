package database

import (
	"fmt"
)

type ConvSnippet struct {
	ConversationID   uint
	ConversationType string // "private" or "group"
	lastMessage      *MessageSnippet
}
type MessageSnippet struct {
	Sender     string
	Content    string
	Timestamp  string
	Attachment []byte
	Status     string
}

func (db *appdbimpl) GetConversationSnippet(conversationID uint) (*ConvSnippet, error) {
	var messageSnippet MessageSnippet
	var snippet ConvSnippet
	snippet.lastMessage = &messageSnippet
	err := db.c.QueryRow(`
		SELECT id, type
		FROM Conversations
		WHERE id = ?
		`, conversationID).Scan(&snippet.ConversationID, &snippet.ConversationType)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation snippet: %w", err)
	}

	if snippet.ConversationType == "private" {
		snippet.ConversationType = "private"
		err = db.c.QueryRow(`
			SELECT m.sender_username, m.content, m.timestamp, m.attachment, m.status
			FROM Conversations c JOIN Messages m ON c.id = m.conversation_id
			WHERE c.id = ?
			ORDER BY m.timestamp DESC
			LIMIT 1
			`, conversationID).Scan(&messageSnippet.Sender, &messageSnippet.Content, &messageSnippet.Timestamp, &messageSnippet.Attachment, &messageSnippet.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to get last message for private conversation: %w", err)
		}
	} else if snippet.ConversationType == "group" {
		snippet.ConversationType = "group"
		err = db.c.QueryRow(`
			SELECT m.sender_username, m.content, m.timestamp, m.attachment, m.status
			FROM Conversations c JOIN Messages m ON c.id = m.conversation_id
			WHERE c.id = ?
			ORDER BY m.timestamp DESC
			LIMIT 1
			`, conversationID).Scan(&messageSnippet.Sender, &messageSnippet.Content, &messageSnippet.Timestamp, &messageSnippet.Attachment, &messageSnippet.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to get last message for group conversation: %w", err)
		}
	}
	return &snippet, nil
}

func (db *appdbimpl) GetAllConIDsForUser(username string) ([]uint, error) {
	rows, err := db.c.Query(`
		SELECT DISTINCT c.id
		FROM Conversations c
		LEFT JOIN Private_conversations_memberships pcm ON c.id = pcm.conversation_id
		LEFT JOIN Groups g ON c.id = g.conversation_id
		LEFT JOIN Groups_memberships gm ON g.id = gm.group_id
		WHERE (c.type = 'private' AND pcm.member_username = ?)
		   OR (c.type = 'group' AND gm.member_username = ?);
		`, username, username)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversation IDs for user: %w", err)
	}
	defer rows.Close()

	var conversationIDs []uint
	for rows.Next() {
		var convID uint
		if err := rows.Scan(&convID); err != nil {
			return nil, fmt.Errorf("failed to scan conversation ID: %w", err)
		}
		conversationIDs = append(conversationIDs, convID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over conversation IDs: %w", err)
	}

	return conversationIDs, nil
}

// GetConversationParticipants retrieves all usernames participating in a conversation
func (db *appdbimpl) GetConversationParticipants(conversationID uint) ([]string, error) {
	// First check the conversation type
	var convType string
	err := db.c.QueryRow(`SELECT type FROM Conversations WHERE id = ?`, conversationID).Scan(&convType)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation type: %w", err)
	}

	var participants []string

	if convType == "private" {
		// Get participants from Private_conversations_memberships
		rows, err := db.c.Query(`
			SELECT member_username
			FROM Private_conversations_memberships
			WHERE conversation_id = ?
		`, conversationID)
		if err != nil {
			return nil, fmt.Errorf("failed to query private conversation participants: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var username string
			if err := rows.Scan(&username); err != nil {
				return nil, fmt.Errorf("failed to scan participant: %w", err)
			}
			participants = append(participants, username)
		}

		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating participants: %w", err)
		}
	} else if convType == "group" {
		// Get participants from Groups_memberships via group ID
		rows, err := db.c.Query(`
			SELECT gm.member_username
			FROM Groups g
			JOIN Groups_memberships gm ON g.id = gm.group_id
			WHERE g.conversation_id = ?
		`, conversationID)
		if err != nil {
			return nil, fmt.Errorf("failed to query group conversation participants: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var username string
			if err := rows.Scan(&username); err != nil {
				return nil, fmt.Errorf("failed to scan participant: %w", err)
			}
			participants = append(participants, username)
		}

		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating participants: %w", err)
		}
	}

	return participants, nil
}
