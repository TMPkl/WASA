package database

import (
	"errors"
	"fmt"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database/attachments"
)

/*
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	conversation_id INTEGER NOT NULL,
	sender_username TEXT NOT NULL,
	content
	 TEXT,
	timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
	attachment BLOB DEFAULT NULL,
	reaction text DEFAULT NULL,
	status TEXT CHECK( status IN ('sent','delivered','received','hidden') ), /// musze sprawdzic czy ja to 100% musze zrobic
	CHECK (content IS NOT NULL OR attachment IS NOT NULL)
*/

type Message struct {
	ID             int64
	ConversationID uint
	SenderUsername string
	Content        string
	Timestamp      time.Time
	Attachment     []byte // GOB encoded AttachmentsPack
	Reaction       string
	Status         string
}

func (db *appdbimpl) SaveMessage(username string, MessageContent string, ap attachments.AttachmentsPack, ConvID uint) (Message, error) {
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
	}
	attachmentsData, err := ap.ConvertToGOB()
	if err != nil {
		return Message{}, fmt.Errorf("Failed to convert attachments to GOB: %w", err)
	}
	message.Attachment = attachmentsData

	res, err := db.c.Exec("INSERT INTO Messages (conversation_id, sender_username, content, timestamp, attachment, status) VALUES (?, ?, ?, ?, ?, ?)",
		message.ConversationID, message.SenderUsername, message.Content, message.Timestamp, message.Attachment, message.Status)
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
