package api

//################file for endpoints from tag conversations##########
import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")
	if username == "" {
		// fallback to query param if path param missing
		username = r.URL.Query().Get("username")
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
	conIDs, err := rt.db.GetAllConIDsForUser(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get conversations: %v", err), http.StatusInternalServerError)
		return
	}

	var conversations []database.ConvSnippet
	for _, conID := range conIDs {
		snippet, err := rt.db.GetConversationSnippet(conID)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get conversation snippet: %v", err), http.StatusInternalServerError)
			return
		}
		conversations = append(conversations, *snippet)
	}
	data, err := json.Marshal(conversations)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
func (rt *_router) GetConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type rqst struct {
		Username        string `json:"username"`
		MessegeQuantity int    `json:"message_quantity"`
	}

	conversationIDStr := ps.ByName("conversationId")
	if conversationIDStr == "" {
		// fallback to query param if path param missing
		conversationIDStr = r.URL.Query().Get("conversationId")
	}
	// Convert conversationIDStr to uint
	var conversationID uint
	_, err := fmt.Sscanf(conversationIDStr, "%d", &conversationID)
	if err != nil {
		http.Error(w, "invalid conversation ID", http.StatusBadRequest)
		return
	}

	var request rqst
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	authorised, err := rt.Authorise(w, r, request.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("authorization error: %v", err), http.StatusUnauthorized)
		return
	}
	if !authorised {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	isInConv, err := rt.db.UserInConversation(request.Username, conversationID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to verify user in conversation: %v", err), http.StatusInternalServerError)
		return
	}
	if !isInConv {
		http.Error(w, "user not in conversation", http.StatusForbidden)
		return
	}

	// Get conversation type
	convType, err := rt.db.GetConversationType(conversationID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get conversation type: %v", err), http.StatusInternalServerError)
		return
	}
	isGroup := convType == "group"

	// Get conversation participants
	participants, err := rt.db.GetConversationParticipants(conversationID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get conversation participants: %v", err), http.StatusInternalServerError)
		return
	}

	// Get message history
	limit := request.MessegeQuantity
	if limit <= 0 || limit > 100 {
		limit = 50 // Default to 50 messages
	}

	messages, err := rt.db.GetConversationHistory(conversationID, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get conversation history: %v", err), http.StatusInternalServerError)
		return
	}

	// Build response
	type MessageResponse struct {
		ID             int64    `json:"id"`
		SenderUsername string   `json:"sender_username"`
		Content        string   `json:"content"`
		Timestamp      string   `json:"timestamp"`
		HasAttachment  bool     `json:"has_attachment"`
		Reactions      []string `json:"reactions"`
		Status         string   `json:"status"`
	}

	type ConversationResponse struct {
		ID           uint              `json:"conversation_id"`
		IsGroup      bool              `json:"is_group"`
		Participants []string          `json:"participants"`
		Messages     []MessageResponse `json:"messages"`
	}

	var messageResponses []MessageResponse
	for _, msg := range messages {
		// Get reactions for this message
		reactionIDs, _ := rt.db.GetReactionsID(int(msg.ID))
		reactions := []string{}
		for _, reactionID := range reactionIDs {
			reactionMsg, err := rt.db.GetMessageByID(fmt.Sprintf("%d", reactionID))
			if err == nil {
				reactions = append(reactions, fmt.Sprintf("%s: %s", reactionMsg.SenderUsername, reactionMsg.Content))
			}
		}

		messageResponses = append(messageResponses, MessageResponse{
			ID:             msg.ID,
			SenderUsername: msg.SenderUsername,
			Content:        msg.Content,
			Timestamp:      msg.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			HasAttachment:  len(msg.Attachment) > 0,
			Reactions:      reactions,
			Status:         msg.Status,
		})
	}

	response := ConversationResponse{
		ID:           conversationID,
		IsGroup:      isGroup,
		Participants: participants,
		Messages:     messageResponses,
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
