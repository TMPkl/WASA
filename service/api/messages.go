package api

//################file for endpoints form tag messages##########
import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database/attachments"
	"github.com/julienschmidt/httprouter"
)

type SendMessageRequest struct {
	Content         string
	ConversationID  int
	SenderUsername  string
	ReciverUsername string
	ReplyingToID    int64
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

	conversationIdStr := r.FormValue("conversationId")
	if conversationIdStr != "" {
		_, err := fmt.Sscanf(conversationIdStr, "%d", &rqst.ConversationID)
		if err != nil {
			rt.baseLogger.Printf("Invalid conversationId: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	replyingToIdStr := r.FormValue("replyingToId")
	if replyingToIdStr != "" {
		_, err := fmt.Sscanf(replyingToIdStr, "%d", &rqst.ReplyingToID)
		if err != nil {
			rt.baseLogger.Printf("Invalid replyingToId: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if rqst.Content == "" && len(r.MultipartForm.File["attachments"]) == 0 {
		rt.baseLogger.Printf("Message content is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var czyJestGrupa bool
	if rqst.ConversationID == 0 {
		// No conversation ID provided — this is a new/private conversation flow
		czyJestGrupa = false
	} else {
		var err error
		czyJestGrupa, err = rt.db.IsConversationGroup(uint(rqst.ConversationID))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				rt.baseLogger.Printf("Conversation not found: %s", err.Error())
				w.WriteHeader(http.StatusNotFound)
				return
			}
			rt.baseLogger.Printf("Failed to check if conversation is group: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if rqst.ReciverUsername == "" && !czyJestGrupa {
		rt.baseLogger.Printf("Receiver username is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if czyJestGrupa {
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

		var message database.Message
		var err error
		if rqst.ReplyingToID > 0 {
			message, err = rt.db.SaveMessageWithReply(rqst.SenderUsername, rqst.Content, attachmentsPack, uint(rqst.ConversationID), int64(rqst.ReplyingToID))
		} else {
			message, err = rt.db.SaveMessage(rqst.SenderUsername, rqst.Content, attachmentsPack, uint(rqst.ConversationID))
		}
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
			"conversationId": fmt.Sprintf("%d", rqst.ConversationID),
			"timestamp":      message.Timestamp.Format("2022-01-02T15:04:05Z07:00"),
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			rt.baseLogger.Printf("Failed to encode response: %s", err.Error())
		}

	}

	// check if conversation exists
	convID, err := rt.db.DoesUsersOwnConversation(rqst.SenderUsername, rqst.ReciverUsername)
	if err != nil {
		convID, err = rt.db.CreatePrivateConversation(rqst.SenderUsername, rqst.ReciverUsername)
		if err != nil {
			rt.baseLogger.Printf("Failed to create conversation: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

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

	var message database.Message
	if rqst.ReplyingToID > 0 {
		message, err = rt.db.SaveMessageWithReply(rqst.SenderUsername, rqst.Content, attachmentsPack, convID, int64(rqst.ReplyingToID))
	} else {
		message, err = rt.db.SaveMessage(rqst.SenderUsername, rqst.Content, attachmentsPack, convID)
	}
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
	NewMessage.SenderUsername = rqst.Username

	if len(forwardedMessage.Attachment) > 0 {
		apPtr, err := attachments.DecodeFromGOB(forwardedMessage.Attachment)
		if err != nil {
			rt.baseLogger.Printf("Failed to decode attachments: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if apPtr != nil {
			NewMessage.Attachmemnts = *apPtr
		}
	}

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
	_, _ = fmt.Sscanf(s, "%d", &n)
	return n
}

func (rt *_router) ReactToMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	reactedToMessagID := ps.ByName("messageId")

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

type RemoveReactionRequest struct {
	Username string
}

func (rt *_router) RemoveReactionFromMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	baseMessageID := ps.ByName("messageId")

	var rqst RemoveReactionRequest

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

	conversationID, _ := rt.db.ConversationIDfromMessageID(baseMessageID)
	IsUserInConv, err := rt.db.UserInConversation(rqst.Username, conversationID)
	if err != nil {
		rt.baseLogger.Printf("Failed to check if user is in conversation: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !IsUserInConv {
		rt.baseLogger.Printf("User %s is not in conversation %d", rqst.Username, conversationID)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	reactionMessageID, err := rt.db.GetReactionIDByUsernameAndBaseMessageID(rqst.Username, atoi(baseMessageID))
	if err != nil {
		rt.baseLogger.Printf("Failed : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.db.DeleteMessage(fmt.Sprintf("%d", reactionMessageID))
	if err != nil {
		rt.baseLogger.Printf("Failed to delete reaction message: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) GetMessageAttachments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	messageID := ps.ByName("messageId")
	if messageID == "" {
		http.Error(w, "message ID is required", http.StatusBadRequest)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	authorised, err := rt.Authorise(w, r, username)
	if err != nil {
		http.Error(w, fmt.Sprintf("authorization error: %v", err), http.StatusUnauthorized)
		return
	}
	if !authorised {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	message, err := rt.db.GetMessageByID(messageID)
	if err != nil {
		rt.baseLogger.Printf("Failed to get message by ID %s: %v", messageID, err)
		http.Error(w, "message not found", http.StatusNotFound)
		return
	}

	isInConv, err := rt.db.UserInConversation(username, message.ConversationID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to verify user in conversation: %v", err), http.StatusInternalServerError)
		return
	}
	if !isInConv {
		http.Error(w, "forbidden: user not in conversation", http.StatusForbidden)
		return
	}

	if len(message.Attachment) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"attachments":[]}`))
		return
	}

	attachmentsPack, err := attachments.DecodeFromGOB(message.Attachment)
	if err != nil {
		http.Error(w, "failed to decode attachments", http.StatusInternalServerError)
		return
	}

	type AttachmentResponse struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	}

	type Response struct {
		Attachments []AttachmentResponse `json:"attachments"`
	}

	var attachmentResponses []AttachmentResponse
	for _, att := range attachmentsPack.Attachments {
		attachmentResponses = append(attachmentResponses, AttachmentResponse{
			Type:    att.Type,
			Content: base64.StdEncoding.EncodeToString(att.Content),
		})
	}

	response := Response{
		Attachments: attachmentResponses,
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
func (rt *_router) GetAttachmentFromMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	messageID := ps.ByName("messageId")
	if messageID == "" {
		http.Error(w, "message ID is required", http.StatusBadRequest)
		return
	}

	attachmentData, err := rt.db.GetAttachmentFromMessage(atoi(messageID))
	if err != nil {
		http.Error(w, "failed to get attachment from message", http.StatusInternalServerError)
		return
	}

	if len(attachmentData) == 0 {
		http.Error(w, "no attachments found", http.StatusNotFound)
		return
	}
	// decodie
	attachmentsPack, err := attachments.DecodeFromGOB(attachmentData)
	if err != nil {
		rt.baseLogger.Printf("Failed to decode attachments: %v", err)
		http.Error(w, "failed to decode attachments", http.StatusInternalServerError)
		return
	}

	if attachmentsPack.IsEmpty() {
		http.Error(w, "no attachments found", http.StatusNotFound)
		return
	}

	att := attachmentsPack.Attachments[0]
	w.Header().Set("Content-Type", att.Type)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"attachment_%s\"", messageID))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(att.Content)
}
