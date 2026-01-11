package database

import (
	"fmt"
)

type ConvSnippet struct {
	ConversationID   uint   `json:"ConversationID"`
	ConversationType string `json:"ConversationType"` // "private" or "group"
	GroupName        string `json:"GroupName,omitempty"`
	OtherUsername    string `json:"OtherUsername,omitempty"`
	LastMessage      string `json:"LastMessage"`
	LastMessageTime  string `json:"LastMessageTime"`
	Status           string `json:"Status,omitempty"`
}
type MessageSnippet struct {
	Sender     string
	Content    string
	Timestamp  string
	Attachment []byte
	Status     string
}

func (db *appdbimpl) GetConversationSnippet(conversationID uint) (*ConvSnippet, error) {
	var snippet ConvSnippet
	err := db.c.QueryRow(`
		SELECT id, type
		FROM Conversations
		WHERE id = ?
		`, conversationID).Scan(&snippet.ConversationID, &snippet.ConversationType)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation snippet: %w", err)
	}

	if snippet.ConversationType == "private" {
		// Get the other user's username and last message
		err = db.c.QueryRow(`
			SELECT 
				CASE 
					WHEN pcm1.member_username = pcm2.member_username THEN pcm1.member_username
					ELSE CASE 
						WHEN pcm1.member_username < pcm2.member_username THEN pcm2.member_username
						ELSE pcm1.member_username
					END
				END as other_username,
				COALESCE(m.content, '') as last_message,
				COALESCE(m.timestamp, '') as last_timestamp,
				COALESCE(m.status, '') as last_status
			FROM Private_conversations_memberships pcm1
			LEFT JOIN Private_conversations_memberships pcm2 
				ON pcm1.conversation_id = pcm2.conversation_id 
				AND pcm1.member_username != pcm2.member_username
			LEFT JOIN (
				SELECT conversation_id, content, timestamp, status
				FROM Messages
				WHERE conversation_id = ?
				ORDER BY timestamp DESC
				LIMIT 1
			) m ON pcm1.conversation_id = m.conversation_id
			WHERE pcm1.conversation_id = ?
			LIMIT 1
		`, conversationID, conversationID).Scan(&snippet.OtherUsername, &snippet.LastMessage, &snippet.LastMessageTime, &snippet.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to get private conversation details: %w", err)
		}
	} else if snippet.ConversationType == "group" {
		// Get group name and last message
		err = db.c.QueryRow(`
			SELECT 
				g.name,
				COALESCE(m.content, '') as last_message,
				COALESCE(m.timestamp, '') as last_timestamp,
				COALESCE(m.status, '') as last_status
			FROM Groups g
			LEFT JOIN (
				SELECT conversation_id, content, timestamp, status
				FROM Messages
				WHERE conversation_id = ?
				ORDER BY timestamp DESC
				LIMIT 1
			) m ON g.conversation_id = m.conversation_id
			WHERE g.conversation_id = ?
		`, conversationID, conversationID).Scan(&snippet.GroupName, &snippet.LastMessage, &snippet.LastMessageTime, &snippet.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to get group conversation details: %w", err)
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
