package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Profile endpoints
	rt.router.POST("/login", rt.login)
	rt.router.PATCH("/me/name", rt.UpdateMyUsername)
	rt.router.POST("/me/photo", rt.SetProfilePhoto)

	//Messages endpoints

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
