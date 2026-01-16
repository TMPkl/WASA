package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Profile endpoints
	rt.router.POST("/login", rt.login)
	rt.router.PATCH("/me/username", rt.UpdateMyUsername)
	rt.router.POST("/me/photo", rt.SetProfilePhoto)

	rt.router.GET("/users/:username/photo", rt.GetUserProfilePhoto)

	//Messages endpoints
	rt.router.POST("/messages", rt.SendMessage)
	rt.router.DELETE("/messages/:messageId", rt.DeleteMessage)
	rt.router.PATCH("/messages/:messageId/status", rt.UpdateMessageStatus)
	// Forward message
	rt.router.POST("/messages/:messageId/forwards", rt.ForwardMessage)

	rt.router.POST("/messages/:messageId/reactions", rt.ReactToMessage)
	rt.router.DELETE("/messages/:messageId/reactions", rt.RemoveReactionFromMessage)
	rt.router.GET("/messages/:messageId/attachments", rt.GetMessageAttachments)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// Groups endpoints
	rt.router.POST("/groups", rt.CreateGroup)
	rt.router.POST("/groups/:groupId/members", rt.AddGroupMember)
	rt.router.DELETE("/groups/:groupId/members/me", rt.RemoveMeFromGroup)

	rt.router.PATCH("/groups/:groupId/name", rt.RenameGroup)
	rt.router.POST("/groups/:groupId/photo", rt.UpdateGroupPhoto)
	rt.router.GET("/groups/:groupId/photo", rt.GetGroupPhoto)

	rt.router.GET("/conversations/:username", rt.GetMyConversations)
	rt.router.POST("/conversation/:conversationId", rt.GetConversation)
	rt.router.POST("/attachments/:messageId", rt.GetAttachmentFromMessage)

	rt.router.GET("/users", rt.SearchUsers)
	rt.router.POST("/conversations/list", rt.GetAllMyConversations)
	//test
	// rt.router.POST("/test", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// 	data, err := rt.readFileFromRequest(r)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write(data)
	// })
	return rt.router
}
