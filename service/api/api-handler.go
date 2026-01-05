package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Test endpoint
	//rt.router.GET("/test", rt.checkUser)

	// Profile endpoints
	rt.router.POST("/login", rt.login)
	rt.router.PATCH("/me/name", rt.UpdateMyUsername)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
