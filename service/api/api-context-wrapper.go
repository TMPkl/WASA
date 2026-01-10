package api

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

// wrapBase parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrapBase(w http.ResponseWriter, r *http.Request) (*reqcontext.RequestContext, bool) {
	reqUUID, err := uuid.NewV4()
	if err != nil {
		rt.baseLogger.WithError(err).Error("can't generate a request UUID")
		w.WriteHeader(http.StatusInternalServerError)
		return nil, false
	}
	ctx := &reqcontext.RequestContext{
		ReqUUID: reqUUID,
	}

	// Create a request-specific logger
	ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
		"reqid":     ctx.ReqUUID.String(),
		"remote-ip": r.RemoteAddr,
	})

	return ctx, true
}
