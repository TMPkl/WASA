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
