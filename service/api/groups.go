package api

//################file for endpoints from tag groups##########
import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) CreateGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type request struct {
		GroupName string   `json:"group_name"`
		Members   []string `json:"members"`
		Username  string   `json:"username"` // creator username
	}
	type response struct {
		GroupID uint `json:"group_id"`
	}

	var rqst request
	err := json.NewDecoder(r.Body).Decode(&rqst)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if rqst.GroupName == "" {
		http.Error(w, "group_name is required", http.StatusBadRequest)
		return
	}

	authorised, err := rt.Authorise(w, r, rqst.Username)

	if err != nil {
		http.Error(w, fmt.Sprintf("authorization error: %v", err), http.StatusUnauthorized)
		return
	}
	if !authorised {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := rt.db.CreateGroup(rqst.GroupName, rqst.Members)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create group: %v", err), http.StatusInternalServerError)
		return
	}

	resp := response{GroupID: groupID}
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(data)
}
