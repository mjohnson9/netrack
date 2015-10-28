package app

import (
	"fmt"
	"net/http"

	"github.com/mjibson/appstats"

	"appengine"
)

func CreateHandler(handlerFunc func(appengine.Context, http.ResponseWriter, *http.Request)) appstats.Handler {
	return appstats.NewHandler(func(c appengine.Context, w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Set("X-Frame-Options", "SAMEORIGIN")
		headers.Set("X-XSS-Protection", "0")

		c.Infof("Request ID: %s", appengine.RequestID(c))

		handlerFunc(c, w, r)
	})
}

func ServerError(w http.ResponseWriter, r *http.Request) {
	SetCacheControl(w, CacheNoStore, -1)

	c := appengine.NewContext(r)
	requestID := appengine.RequestID(c)
	http.Error(w, fmt.Sprintf("There was an internal server error while processing your request. Administrators have been notified.\nRequest ID: %s", requestID), http.StatusInternalServerError)
}
