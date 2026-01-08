package api

//################file for endpoints form tag messages##########
import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database/attachments"
	"github.com/julienschmidt/httprouter"
)

type SendMessageRequest struct {
	Content              string //`json:"content"`
	RespondedToMessageID int    //`json:"respondedToMessageId"`
	ConversationID       int    //`json:"conversationId"`/// trzeba tez ustalic jak wygladaje te id
	SenderUsername       string //`json:"senderUsername"`
	ReciverUsername      string //`json:"receiverUsername"`
	Attachmemnts         attachments.AttachmentsPack
}

func (rt *_router) SendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var rqst SendMessageRequest

	err := r.ParseMultipartForm(10 << 22)
	if err != nil {
		rt.baseLogger.Printf("Error parsing multipart form")
		rt.baseLogger.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rqst.SenderUsername = r.FormValue("senderUsername")
	if rqst.SenderUsername == "" {
		rt.baseLogger.Printf("Username is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authorised, err := rt.Authorise(w, r, rqst.SenderUsername)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !authorised {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rqst.Content = r.FormValue("content")
	rqst.ReciverUsername = r.FormValue("receiverUsername")

	if rqst.Content == "" {
		rt.baseLogger.Printf("Message content is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if rqst.ReciverUsername == "" {
		rt.baseLogger.Printf("Receiver username is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if conversation exists
	convID, err := rt.db.DoesUsersOwnConversation(rqst.SenderUsername, rqst.ReciverUsername)
	if err != nil {
		// Conversation doesn't exist, create new one
		rt.baseLogger.Printf("Creating new conversation between %s and %s", rqst.SenderUsername, rqst.ReciverUsername)
		convID, err = rt.db.CreatePrivateConversation(rqst.SenderUsername, rqst.ReciverUsername)
		if err != nil {
			rt.baseLogger.Printf("Failed to create conversation: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Parse attachments from multipart form
	var attachmentsPack attachments.AttachmentsPack
	files := r.MultipartForm.File["attachments"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			rt.baseLogger.Printf("Failed to open attachment: %s", err.Error())
			continue
		}
		defer file.Close()

		content := make([]byte, fileHeader.Size)
		_, err = file.Read(content)
		if err != nil {
			rt.baseLogger.Printf("Failed to read attachment: %s", err.Error())
			continue
		}

		attachment := attachments.NewAttachment(content)
		attachmentsPack.Attachments = append(attachmentsPack.Attachments, attachment)
	}

	// Save message to database
	message, err := rt.db.SaveMessage(rqst.SenderUsername, rqst.Content, attachmentsPack, convID)
	if err != nil {
		rt.baseLogger.Printf("Failed to save message: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"id":        fmt.Sprintf("%d", message.ID),
		"senderId":  message.SenderUsername,
		"content":   message.Content,
		"timestamp": message.Timestamp.Format("2022-01-02T15:04:05Z07:00"),
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		rt.baseLogger.Printf("Failed to encode response: %s", err.Error())
	}
}

func (rt *_router) DeleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	/*
	   doc:
	   	/messages/{messageId}:
	   		delete:
	   		tags: [messages]
	   		summary: Delete message
	   		description: Deletes a message.
	   		operationId: deleteMessage
	   		security: [{ bearerAuth: [] }]
	   		parameters:
	   			- $ref: '#/components/parameters/MessageId'
	   		responses:
	   			"204":
	   			description: Message deleted
	*/
	messageId := ps.ByName("messageId")
	_ = messageId

	msg, err := rt.db.GetMessageByID(messageId)

	if err != nil {
		rt.baseLogger.Printf("Failed to get message: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if msg == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	authorised, err := rt.Authorise(w, r, msg.SenderUsername)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !authorised {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = rt.db.DeleteMessage(messageId)
	if err != nil {
		rt.baseLogger.Printf("Failed to delete message: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
