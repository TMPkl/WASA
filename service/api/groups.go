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
		GroupName string `json:"group_name"`
		Username  string `json:"username"` // creator username
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

	groupID, err := rt.db.CreateGroup(rqst.GroupName, rqst.Username)
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

func (rt *_router) AddGroupMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	/*Content type
	application/json

	Copy
	{
	"userId": "^a$"
	}

	https://tmpkl.github.io/groups/{groupId}/members

	*/
	type request struct {
		Username      string `json:"username"`
		UsernemeToAdd string `json:"user_to_add"`
	}
	groupIDParam := ps.ByName("groupId")
	var rqst request

	err := json.NewDecoder(r.Body).Decode(&rqst)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if rqst.Username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	DoesUserEx, err := rt.db.UserExists(rqst.UsernemeToAdd)

	if err != nil {
		http.Error(w, fmt.Sprintf("error checking if user exists: %v", err), http.StatusInternalServerError)
		return
	}
	if DoesUserEx != true {
		http.Error(w, "user to add does not exist", http.StatusBadRequest)
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
	//trzeba sprawdzic czy username jest w grupie bo tylko twdy moze dodawc ludzi chyba tka bylo w dokumentacji?????? alo i ne?

	legit, err := rt.db.IsUserInGroup(uint(atoi(groupIDParam)), rqst.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("error checking group membership: %v", err), http.StatusInternalServerError)
		return
	}
	if !legit {
		http.Error(w, "user not in group", http.StatusUnauthorized)
		return
	}
	// trzena tez sprawdzic czy user nie jest juz w griupie
	userAlreadyInGroup, err := rt.db.IsUserInGroup(uint(atoi(groupIDParam)), rqst.UsernemeToAdd)
	if err != nil {
		http.Error(w, fmt.Sprintf("error checking if user is already in group: %v", err), http.StatusInternalServerError)
		return
	}
	if userAlreadyInGroup {
		http.Error(w, "user already in group", http.StatusAlreadyReported)
		return
	}

	///ehhh phrasing
	var groupID uint
	_, err = fmt.Sscanf(groupIDParam, "%d", &groupID)
	if err != nil {
		http.Error(w, "invalid group ID", http.StatusBadRequest)
		return
	}

	err = rt.db.AddMemberToGroup(groupID, rqst.UsernemeToAdd)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to add member to group: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
