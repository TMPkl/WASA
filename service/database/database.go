/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database/attachments"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database/migrations"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	Ping() error

	UserExists(username string) (bool, error)
	AddNewUser(username string) error
	UpdateUsername(oldUsername, newUsername string) error
	AddProfilePhoto(username string, photoData []byte) error
	GetProfilePhoto(username string) ([]byte, error)
	GetAllUsers() ([]string, error)

	DoesUsersOwnConversation(unA string, unB string) (uint, error)
	CreatePrivateConversation(username1 string, username2 string) (uint, error)
	SaveMessage(username string, MessageContent string, ap attachments.AttachmentsPack, ConvID uint) (Message, error)
	SaveMessageWithReply(username string, MessageContent string, ap attachments.AttachmentsPack, ConvID uint, replyingToID int64) (Message, error)

	GetMessageByID(messageID string) (*Message, error)
	DeleteMessage(messageID string) error
	MessageOwner(messageID string) (string, error)
	UserInConversation(username string, conversationID uint) (bool, error)
	ConversationIDfromMessageID(messageID string) (uint, error)

	ReactToMessage(reactedToMessagID int, reaction string, username string, conversationID uint) (uint, error)
	GetReactionIDByUsernameAndBaseMessageID(username string, baseMessageID int) (int, error)
	GetReactionsID(baseMessageID int) ([]int, bool)

	// GROPS DB
	CreateGroup(groupName string, firstMember string) (uint, error)
	GetGroupMembers(groupID uint) ([]string, error)
	RenameGroup(groupID uint, newName string) error
	AddMemberToGroup(groupID uint, username string) error
	RemoveMemberFromGroup(groupID uint, username string) error
	IsUserInGroup(groupID uint, username string) (bool, error)
	SaveGroupMessage(username string, groupID uint, messageContent string, ap attachments.AttachmentsPack) (Message, error)
	UpdateGroupPhoto(groupID uint, photoData []byte) error
	GetGroupPhoto(groupID uint) ([]byte, error)
	GroupExists(groupID uint) (bool, error)

	GetConversationSnippet(conversationID uint) (*ConvSnippet, error)
	GetAllConIDsForUser(username string) ([]uint, error)
	GetConversationHistory(conversationID uint, limit int) ([]Message, error)
	GetConversationParticipants(conversationID uint) ([]string, error)
	IsConversationGroup(conversationID uint) (bool, error)

	GetConversationType(conversationID uint) (string, error)
	GetAttachmentFromMessage(messageID int) ([]byte, error)

	GetGroupIDByConversationID(conversationID uint) (uint, error)
	GetAllConversations(username string) ([]ConvInfo, error)
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	err := migrations.Migrate(db)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
