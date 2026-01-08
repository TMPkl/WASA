package api

//################file for endpoints form tag messages##########
import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database/attachments"
	"github.com/julienschmidt/httprouter"
)

type SendMessageRequest struct {
	Content         string //`json:"content"`
	ConversationID  int    //`json:"conversationId"`/// trzeba tez ustalic jak wygladaje te id
	SenderUsername  string //`json:"senderUsername"`
	ReciverUsername string //`json:"receiverUsername"`
	Attachmemnts    attachments.AttachmentsPack
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
		"id":             fmt.Sprintf("%d", message.ID),
		"senderId":       message.SenderUsername,
		"content":        message.Content,
		"conversationId": fmt.Sprintf("%d", convID),
		"timestamp":      message.Timestamp.Format("2022-01-02T15:04:05Z07:00"),
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

type ForwardMessageRequest struct {
	AddressingConversationID uint   `json:"addressingConversationID"`
	Username                 string `json:"username"`
}

func (rt *_router) ForwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//sprawdz autentykacje
	//sprawdz czy uzytkownik jest w konwersacji
	//zapisz wiadomosc
	//przekaz wiadmosc z zmienionymi wartosciami
	var rqst ForwardMessageRequest
	err := json.NewDecoder(r.Body).Decode(&rqst)
	if err != nil {
		rt.baseLogger.Printf("Error decoding message request: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authorised, err := rt.Authorise(w, r, rqst.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !authorised {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	inConv, err := rt.db.UserInConversation(rqst.Username, rqst.AddressingConversationID)
	if err != nil {
		rt.baseLogger.Printf("Failed to wzium wzium: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !inConv {
		rt.baseLogger.Printf("User %s is not in conversation %s", rqst.Username, rqst.AddressingConversationID)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	mesgID := ps.ByName("messageId")
	forwardedMessage, err := rt.db.GetMessageByID(mesgID)
	if err != nil {
		rt.baseLogger.Printf("Failed to get message: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if forwardedMessage == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var NewMessage SendMessageRequest
	NewMessage.Content = forwardedMessage.Content
	NewMessage.ConversationID = int(rqst.AddressingConversationID)
	apPtr, err := attachments.DecodeFromGOB(forwardedMessage.Attachment)
	if err != nil {
		rt.baseLogger.Printf("Failed to decode attachments: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if apPtr != nil {
		NewMessage.Attachmemnts = *apPtr
	}
	NewMessage.SenderUsername = rqst.Username

	_, err = rt.db.SaveMessage(NewMessage.SenderUsername, NewMessage.Content, NewMessage.Attachmemnts, uint(NewMessage.ConversationID))
	if err != nil {
		rt.baseLogger.Printf("Failed to save forwarded message: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

/*type SendMessageRequest struct {
	Content              string //`json:"content"`
	RespondedToMessageID int    //`json:"respondedToMessageId"`
	ConversationID       int    //`json:"conversationId"`/// trzeba tez ustalic jak wygladaje te id
	SenderUsername       string //`json:"senderUsername"`
	ReciverUsername      string //`json:"receiverUsername"`
	Attachmemnts         attachments.AttachmentsPack
}*/

type ReactionRequest struct {
	Emoji    string `json:"emoji"`
	Username string `json:"username"`
}

func atoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

func (rt *_router) ReactToMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	reactedToMessagID := ps.ByName("messageId")
	_ = reactedToMessagID

	var rqst ReactionRequest
	err := json.NewDecoder(r.Body).Decode(&rqst) // decoder json
	if err != nil {
		rt.baseLogger.Printf("Error decoding KSON: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authorised, err := rt.Authorise(w, r, rqst.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !authorised {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//sprawdz czy uzytkownik jest wkonwersacji z wiadomscia
	//sprawdz czy wiadomosc istnieje
	//zapisz reakcje
	ConversationIDfromRQ, err := rt.db.ConversationIDfromMessageID(reactedToMessagID)
	_ = ConversationIDfromRQ // test

	if err != nil {
		rt.baseLogger.Printf("Failed to get conversation ID from message ID: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	IsUserInConv, err := rt.db.UserInConversation(rqst.Username, ConversationIDfromRQ)
	if err != nil {
		rt.baseLogger.Printf("Failed to check if user is in conversation: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !IsUserInConv {
		rt.baseLogger.Printf("User %s is not in conversation %d", rqst.Username, ConversationIDfromRQ)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	newMessageID, err := rt.db.ReactToMessage(atoi(reactedToMessagID), rqst.Emoji, rqst.Username, ConversationIDfromRQ)
	if err != nil {
		rt.baseLogger.Printf("Failed to save reaction message: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"id":             fmt.Sprintf("%d", newMessageID),
		"senderId":       rqst.Username,
		"content":        rqst.Emoji,
		"conversationId": fmt.Sprintf("%d", ConversationIDfromRQ),
		"timestamp":      time.Now(),
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		rt.baseLogger.Printf("Failed to encode response: %s", err.Error())
	}

}
