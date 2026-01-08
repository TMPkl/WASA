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

	//Messages endpoints
	rt.router.POST("/messages", rt.SendMessage)
	rt.router.DELETE("/messages/:messageId", rt.DeleteMessage)
	// Forward message
	rt.router.POST("/messages/:messageId/forwards", rt.ForwardMessage)

	rt.router.POST("/messages/:messageId/reactions", rt.ReactToMessage)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

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
